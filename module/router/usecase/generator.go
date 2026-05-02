package usecase

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/kodefluence/altair/entity"
	"github.com/kodefluence/altair/module"
)

// defaultUpstreamTimeout caps how long a single upstream request may take.
// 30s mirrors the conservative default for outbound HTTP in services that
// don't have a more specific SLA. Override per-deployment via
// WithUpstreamTimeout (and once exposed, via app.yml).
const defaultUpstreamTimeout = 30 * time.Second

type generatorConfig struct {
	upstreamTimeout    time.Duration
	upstreamTransport  http.RoundTripper
	proxyHost          string
	maxRequestBodySize int64
}

// Option configures a Generator. Add a knob via WithXxx; do not extend the
// positional NewGenerator signature.
type Option func(*generatorConfig)

// WithUpstreamTimeout bounds the entire upstream round-trip (dial + send +
// read body). A timeout of 0 disables it; prefer an explicit upper bound in
// production — a hung upstream otherwise leaks goroutines and FDs forever.
func WithUpstreamTimeout(d time.Duration) Option {
	return func(c *generatorConfig) { c.upstreamTimeout = d }
}

// WithUpstreamTransport overrides the http.RoundTripper used to call
// upstreams. Tests inject this to assert client behavior without standing
// up a real TCP server; production code should rely on the default shared
// transport.
func WithUpstreamTransport(rt http.RoundTripper) Option {
	return func(c *generatorConfig) { c.upstreamTransport = rt }
}

// WithProxyHost sets the Host header sent on every outbound proxy request.
// Captured once at Generator construction so we don't pay an os.Getenv
// syscall per request, and so the value is the same one the rest of the
// app sees via core.AppConfig.ProxyHost().
func WithProxyHost(host string) Option {
	return func(c *generatorConfig) { c.proxyHost = host }
}

// WithMaxRequestBodySize caps inbound request body bytes. Anything larger
// short-circuits with 413. A value <= 0 disables the cap (the historical
// behavior). Use to protect upstreams from clients sending arbitrarily
// large payloads.
func WithMaxRequestBodySize(n int64) Option {
	return func(c *generatorConfig) { c.maxRequestBodySize = n }
}

type Generator struct {
	routerPath         map[string]module.RouterPath
	downStreamPlugin   []module.DownstreamController
	metrics            []module.MetricController
	client             *http.Client
	proxyHost          string
	maxRequestBodySize int64
}

func NewGenerator(downStreamPlugin []module.DownstreamController, metric []module.MetricController, opts ...Option) *Generator {
	cfg := generatorConfig{upstreamTimeout: defaultUpstreamTimeout}
	for _, o := range opts {
		o(&cfg)
	}

	transport := cfg.upstreamTransport
	if transport == nil {
		// One Transport per Generator so connection pooling kicks in across
		// requests. http.DefaultTransport already configures sensible idle/
		// dial timeouts; clone so callers tweaking one Generator don't
		// mutate the package-level default.
		transport = http.DefaultTransport.(*http.Transport).Clone()
	}

	return &Generator{
		routerPath:         map[string]module.RouterPath{},
		downStreamPlugin:   downStreamPlugin,
		metrics:            metric,
		proxyHost:          cfg.proxyHost,
		maxRequestBodySize: cfg.maxRequestBodySize,
		client: &http.Client{
			Transport: transport,
			Timeout:   cfg.upstreamTimeout,
		},
	}
}

func (g *Generator) Generate(engine *gin.Engine, routeObjects []entity.RouteObject) (errVariable error) {
	for _, m := range g.metrics {
		m.InjectCounter("routes_downstream_hits", "route_name", "method", "path", "status_code", "status_code_group")
		m.InjectHistogram("routes_downstream_latency_seconds", "route_name", "method", "path", "status_code", "status_code_group")
		m.InjectHistogram("routes_downstream_plugin_latency_seconds", "route_name", "plugin_name", "method", "path", "status_code", "status_code_group")
	}

	defer func() {
		if r := recover(); r != nil {
			errVariable = fmt.Errorf("Error generating route because of %v", r)
			log.Error().Err(fmt.Errorf("Error generating route because of %v", r)).Array("tags", zerolog.Arr().Str("route").Str("generator").Str("defer").Str("panic")).Msg("Panic error when generating routes")
		}
	}()

	for _, routeObject := range routeObjects {
		for r, routePath := range routeObject.Path {
			g.inheritRouterObject(routeObject, &routePath)

			urlPath := fmt.Sprintf("%s%s", routeObject.Prefix, r)

			g.routerPath[urlPath] = routePath

			log.Info().Str("host", routeObject.Host).Str("name", routeObject.Name).Str("path", urlPath).Array("tags", zerolog.Arr().Str("route").Str("generator").Str("generate").Str("url_path")).Msg("Generating routes")

			engine.Any(urlPath, func(c *gin.Context) {
				requestID := uuid.New().String()
				startTime := time.Now()

				g.do(c, urlPath, requestID, routeObject)

				log.Info().Str("request_id", requestID).Str("host", routeObject.Host).Str("prefix", routeObject.Prefix).Str("name", routeObject.Name).Str("path", urlPath).Str("method", c.Request.Method).Str("full_path", c.Request.URL.String()).Str("client_ip", c.ClientIP()).Float64("duration_seconds", time.Since(startTime).Seconds()).Array("tags", zerolog.Arr().Str("route").Str("generator").Str("generate")).Msg("Complete forwarding the request")
			})
		}
	}

	return errVariable
}

func (g *Generator) do(c *gin.Context, urlPath, requestID string, routeObject entity.RouteObject) {
	proxyReq, err := g.decorateProxyRequest(c, urlPath, requestID, routeObject)
	if err != nil {
		return
	}

	g.decorateHeader(c, requestID, proxyReq)

	if err := g.downStreamPluginCallback(c, proxyReq, urlPath, requestID, routeObject); err != nil {
		return
	}

	if err := g.callDownStreamService(c, proxyReq, urlPath, requestID, routeObject); err != nil {
		return
	}

}

func (g *Generator) decorateProxyRequest(c *gin.Context, urlPath, requestID string, routeObject entity.RouteObject) (*http.Request, error) {
	var proxyReq *http.Request

	if c.Request.Body != nil {
		// MaxBytesReader returns *http.MaxBytesError on overflow. Treat as
		// 413 to give the client an actionable status; ReadAll's other
		// errors still yield 400 (malformed). The cap is opt-in via
		// WithMaxRequestBodySize — zero means "unbounded" to preserve the
		// historical behavior for deployments that haven't set the field.
		reader := c.Request.Body
		if g.maxRequestBodySize > 0 {
			reader = http.MaxBytesReader(c.Writer, c.Request.Body, g.maxRequestBodySize)
		}

		body, err := io.ReadAll(reader)
		if err != nil {
			var maxBytesErr *http.MaxBytesError
			if errors.As(err, &maxBytesErr) {
				log.Warn().Err(err).Stack().Str("host", routeObject.Host).Str("request_id", requestID).Str("prefix", routeObject.Prefix).Str("name", routeObject.Name).Str("path", urlPath).Str("method", c.Request.Method).Str("client_ip", c.ClientIP()).Int64("limit_bytes", g.maxRequestBodySize).Array("tags", zerolog.Arr().Str("route").Str("generator").Str("generate").Str("body_too_large")).Msg("Request body exceeded configured cap")
				c.JSON(http.StatusRequestEntityTooLarge, gin.H{
					"status":  http.StatusRequestEntityTooLarge,
					"message": "Request body too large",
				})
				return nil, err
			}
			log.Error().Err(err).Stack().Str("host", routeObject.Host).Str("request_id", requestID).Str("prefix", routeObject.Prefix).Str("name", routeObject.Name).Str("path", urlPath).Str("method", c.Request.Method).Str("full_path", c.Request.URL.String()).Str("client_ip", c.ClientIP()).Array("tags", zerolog.Arr().Str("route").Str("generator").Str("generate").Str("read_all_request")).Msg("Error reading incoming request body")
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Malformed request body given by the client",
			})
			return nil, err
		}

		proxyReq, err = http.NewRequest(c.Request.Method, "", bytes.NewReader(body))
		if err != nil {
			log.Error().Err(err).Stack().Str("host", routeObject.Host).Str("request_id", requestID).Str("prefix", routeObject.Prefix).Str("name", routeObject.Name).Str("path", urlPath).Str("method", c.Request.Method).Str("full_path", c.Request.URL.String()).Str("client_ip", c.ClientIP()).Array("tags", zerolog.Arr().Str("route").Str("generator").Str("generate").Str("new_request")).Msg("Error creating proxy request")
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Malformed request body given by the client",
			})
			return nil, err
		}
	} else {
		var err error

		proxyReq, err = http.NewRequest(c.Request.Method, "", nil)
		if err != nil {
			log.Error().Err(err).Stack().Str("host", routeObject.Host).Str("request_id", requestID).Str("prefix", routeObject.Prefix).Str("name", routeObject.Name).Str("path", urlPath).Str("method", c.Request.Method).Str("full_path", c.Request.URL.String()).Str("client_ip", c.ClientIP()).Array("tags", zerolog.Arr().Str("route").Str("generator").Str("generate").Str("new_request")).Msg("Error creating proxy request")
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Malformed request body given by the client",
			})
			return nil, err
		}
	}

	proxyReq.URL.Scheme = "http"
	proxyReq.URL.Host = routeObject.Host
	proxyReq.URL.Path = c.Request.URL.Path
	proxyReq.URL.RawQuery = c.Request.URL.RawQuery

	return proxyReq, nil
}

func (g *Generator) decorateHeader(c *gin.Context, requestID string, proxyReq *http.Request) {
	for header, values := range c.Request.Header {
		for _, value := range values {
			proxyReq.Header.Add(header, value)
		}
	}

	proxyReq.Host = g.proxyHost
	proxyReq.Header.Add("X-Request-ID", requestID)
	proxyReq.Header.Set("X-Real-Ip-Address", c.ClientIP())
	proxyReq.Header.Set("X-Forwarded-For", c.Request.RemoteAddr)
}

func (g *Generator) downStreamPluginCallback(c *gin.Context, proxyReq *http.Request, urlPath, requestID string, routeObject entity.RouteObject) error {
	for _, plugin := range g.downStreamPlugin {
		startTimePlugin := time.Now()
		if err := plugin.Intervene(c, proxyReq, g.routerPath[urlPath]); err != nil {
			log.Error().Err(err).Stack().Str("host", routeObject.Host).Str("request_id", requestID).Str("prefix", routeObject.Prefix).Str("name", routeObject.Name).Str("path", urlPath).Str("method", c.Request.Method).Str("full_path", c.Request.URL.String()).Str("client_ip", c.ClientIP()).Array("tags", zerolog.Arr().Str("route").Str("generator").Str("generate").Str("plugin").Str(plugin.Name())).Msg("Plugin error")
			g.downStreamPluginMetric(c, routeObject.Name, plugin.Name(), urlPath, startTimePlugin)
			return err
		}

		log.Info().
			Str("plugin_name", plugin.Name()).
			Float64("plugin_duration_seconds", time.Since(startTimePlugin).Seconds()).
			Array("tags", zerolog.Arr().Str("route").Str("generator").Str("generate").Str("plugin")).
			Msg("Plugin success")

		g.downStreamPluginMetric(c, routeObject.Name, plugin.Name(), urlPath, startTimePlugin)
	}

	return nil
}

func (g *Generator) callDownStreamService(c *gin.Context, proxyReq *http.Request, urlPath, requestID string, routeObject entity.RouteObject) error {
	defer g.downStreamMetric(c, routeObject.Name, urlPath, time.Now())

	proxyRes, err := g.client.Do(proxyReq)
	if err != nil {
		log.Error().Err(err).Stack().Str("host", routeObject.Host).Str("request_id", requestID).Str("prefix", routeObject.Prefix).Str("name", routeObject.Name).Str("path", urlPath).Str("method", c.Request.Method).Str("full_path", c.Request.URL.String()).Str("client_ip", c.ClientIP()).Array("tags", zerolog.Arr().Str("route").Str("generator").Str("generate").Str("client_do")).Msg("Error fowarding the request")
		c.JSON(http.StatusBadGateway, gin.H{
			"status":  http.StatusBadGateway,
			"message": "Bad gateway",
		})
		return err
	}
	defer proxyRes.Body.Close()

	for header, values := range proxyRes.Header {
		for _, value := range values {
			c.Writer.Header().Add(header, value)
		}
	}

	c.Status(proxyRes.StatusCode)

	// Stream the upstream body straight into the gin writer rather than
	// io.ReadAll-buffering it first. A 500MB upstream response previously
	// allocated ~500MB of heap per concurrent request — gone now. If the
	// copy fails partway, the status code is already on the wire so we
	// can't overwrite it with a 502; just log and return.
	if _, err := io.Copy(c.Writer, proxyRes.Body); err != nil {
		log.Error().Err(err).Stack().Str("host", routeObject.Host).Str("request_id", requestID).Str("prefix", routeObject.Prefix).Str("name", routeObject.Name).Str("path", urlPath).Str("method", c.Request.Method).Str("full_path", c.Request.URL.String()).Str("client_ip", c.ClientIP()).Array("tags", zerolog.Arr().Str("route").Str("generator").Str("generate").Str("copy_response")).Msg("Error streaming the response")
		return err
	}

	return nil
}

func (g *Generator) downStreamPluginMetric(c *gin.Context, routeName, pluginName, path string, startTime time.Time) {
	labels := map[string]string{
		"route_name":        routeName,
		"plugin_name":       pluginName,
		"method":            c.Request.Method,
		"path":              path,
		"status_code":       strconv.Itoa(c.Writer.Status()),
		"status_code_group": strconv.Itoa(((c.Writer.Status() / 100) * 100)),
	}

	for _, m := range g.metrics {
		_ = m.Observe("routes_downstream_plugin_latency_seconds", float64(time.Since(startTime).Milliseconds()), labels)
	}
}

func (g *Generator) downStreamMetric(c *gin.Context, routeName, path string, startTime time.Time) {
	labels := map[string]string{
		"route_name":        routeName,
		"method":            c.Request.Method,
		"path":              path,
		"status_code":       strconv.Itoa(c.Writer.Status()),
		"status_code_group": strconv.Itoa(((c.Writer.Status() / 100) * 100)),
	}

	for _, m := range g.metrics {
		_ = m.Inc("routes_downstream_hits", labels)
		_ = m.Observe("routes_downstream_latency_seconds", float64(time.Since(startTime).Milliseconds()), labels)
	}
}

func (g *Generator) inheritRouterObject(routeObject entity.RouteObject, routePath *entity.RouterPath) {
	if routePath.Auth == "" {
		routePath.Auth = routeObject.Auth
	}
}

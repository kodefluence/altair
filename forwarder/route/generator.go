package route

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/altair/entity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type generator struct {
	routerPath       map[string]entity.RouterPath
	downStreamPlugin []core.DownStreamPlugin
	metric           core.Metric
}

// Generator create route generator
func Generator() core.RouteGenerator {
	return &generator{}
}

func (g *generator) Generate(engine *gin.Engine, metric core.Metric, routeObjects []entity.RouteObject, downStreamPlugin []core.DownStreamPlugin) (errVariable error) {
	g.downStreamPlugin = downStreamPlugin

	g.metric = metric
	g.metric.InjectCounter("routes_downstream_hits", "route_name", "method", "path", "status_code", "status_code_group")
	g.metric.InjectHistogram("routes_downstream_latency_seconds", "route_name", "method", "path", "status_code", "status_code_group")
	g.metric.InjectHistogram("routes_downstream_plugin_latency_seconds", "route_name", "plugin_name", "method", "path", "status_code", "status_code_group")

	defer func() {
		if r := recover(); r != nil {
			errVariable = fmt.Errorf("Error generating route because of %v", r)
			log.Error().
				Err(fmt.Errorf("Error generating route because of %v", r)).
				Array("tags", zerolog.Arr().Str("route").Str("generator").Str("defer").Str("panic")).
				Msg("Panic error when generating routes")
		}
	}()

	g.routerPath = map[string]entity.RouterPath{}

	for _, routeObject := range routeObjects {
		for r, routePath := range routeObject.Path {
			g.inheritRouterObject(routeObject, &routePath)

			urlPath := fmt.Sprintf("%s%s", routeObject.Prefix, r)

			g.routerPath[urlPath] = routePath

			log.Info().
				Str("host", routeObject.Host).
				Str("name", routeObject.Name).
				Str("path", urlPath).
				Array("tags", zerolog.Arr().Str("route").Str("generator").Str("generate").Str("url_path")).
				Msg("Generating routes")

			engine.Any(urlPath, func(c *gin.Context) {
				requestID := uuid.New().String()
				startTime := time.Now()

				g.do(c, urlPath, requestID, routeObject)

				log.Info().
					Str("request_id", requestID).
					Str("host", routeObject.Host).
					Str("prefix", routeObject.Prefix).
					Str("name", routeObject.Name).
					Str("path", urlPath).
					Str("method", c.Request.Method).
					Str("full_path", c.Request.URL.String()).
					Str("client_ip", c.ClientIP()).
					Float64("duration_seconds", time.Since(startTime).Seconds()).
					Array("tags", zerolog.Arr().Str("route").Str("generator").Str("generate")).
					Msg("Complete forwarding the request")
			})
		}
	}

	return errVariable
}

func (g *generator) do(c *gin.Context, urlPath, requestID string, routeObject entity.RouteObject) {
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

func (g *generator) decorateProxyRequest(c *gin.Context, urlPath, requestID string, routeObject entity.RouteObject) (*http.Request, error) {
	var proxyReq *http.Request

	if c.Request.Body != nil {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Malformed request body given by the client",
			})
			log.Error().
				Err(err).
				Stack().
				Str("host", routeObject.Host).
				Str("request_id", requestID).
				Str("prefix", routeObject.Prefix).
				Str("name", routeObject.Name).
				Str("path", urlPath).
				Str("method", c.Request.Method).
				Str("full_path", c.Request.URL.String()).
				Str("client_ip", c.ClientIP()).
				Array("tags", zerolog.Arr().Str("route").Str("generator").Str("generate").Str("read_all_request")).
				Msg("Error reading incoming request body")

			return nil, err
		}

		proxyReq, err = http.NewRequest(c.Request.Method, "", bytes.NewReader(body))
		if err != nil {
			log.Error().
				Err(err).
				Stack().
				Str("host", routeObject.Host).
				Str("request_id", requestID).
				Str("prefix", routeObject.Prefix).
				Str("name", routeObject.Name).
				Str("path", urlPath).
				Str("method", c.Request.Method).
				Str("full_path", c.Request.URL.String()).
				Str("client_ip", c.ClientIP()).
				Array("tags", zerolog.Arr().Str("route").Str("generator").Str("generate").Str("new_request")).
				Msg("Error creating proxy request")
			return nil, err
		}
	} else {
		var err error

		proxyReq, err = http.NewRequest(c.Request.Method, "", nil)
		if err != nil {
			log.Error().
				Err(err).
				Stack().
				Str("host", routeObject.Host).
				Str("request_id", requestID).
				Str("prefix", routeObject.Prefix).
				Str("name", routeObject.Name).
				Str("path", urlPath).
				Str("method", c.Request.Method).
				Str("full_path", c.Request.URL.String()).
				Str("client_ip", c.ClientIP()).
				Array("tags", zerolog.Arr().Str("route").Str("generator").Str("generate").Str("new_request")).
				Msg("Error creating proxy request")
			return nil, err
		}
	}

	proxyReq.URL.Scheme = "http"
	proxyReq.URL.Host = routeObject.Host
	proxyReq.URL.Path = c.Request.URL.Path
	proxyReq.URL.RawQuery = c.Request.URL.RawQuery

	return proxyReq, nil
}

func (g *generator) decorateHeader(c *gin.Context, requestID string, proxyReq *http.Request) {
	for header, values := range c.Request.Header {
		for _, value := range values {
			proxyReq.Header.Add(header, value)
		}
	}

	proxyReq.Host = os.Getenv("PROXY_HOST")
	proxyReq.Header.Add("X-Request-ID", requestID)
	proxyReq.Header.Set("X-Real-Ip-Address", c.ClientIP())
	proxyReq.Header.Set("X-Forwarded-For", c.Request.RemoteAddr)
}

func (g *generator) downStreamPluginCallback(c *gin.Context, proxyReq *http.Request, urlPath, requestID string, routeObject entity.RouteObject) error {
	for _, plugin := range g.downStreamPlugin {
		startTimePlugin := time.Now()
		if err := plugin.Intervene(c, proxyReq, g.routerPath[urlPath]); err != nil {
			log.Error().
				Err(err).
				Stack().
				Str("host", routeObject.Host).
				Str("request_id", requestID).
				Str("prefix", routeObject.Prefix).
				Str("name", routeObject.Name).
				Str("path", urlPath).
				Str("method", c.Request.Method).
				Str("full_path", c.Request.URL.String()).
				Str("client_ip", c.ClientIP()).
				Array("tags", zerolog.Arr().Str("route").Str("generator").Str("generate").Str("plugin").Str(plugin.Name())).
				Msg("Plugin error")
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

func (g *generator) callDownStreamService(c *gin.Context, proxyReq *http.Request, urlPath, requestID string, routeObject entity.RouteObject) error {
	defer g.downStreamMetric(c, routeObject.Name, urlPath, time.Now())

	client := http.Client{}
	proxyRes, err := client.Do(proxyReq)
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Str("host", routeObject.Host).
			Str("request_id", requestID).
			Str("prefix", routeObject.Prefix).
			Str("name", routeObject.Name).
			Str("path", urlPath).
			Str("method", c.Request.Method).
			Str("full_path", c.Request.URL.String()).
			Str("client_ip", c.ClientIP()).
			Array("tags", zerolog.Arr().Str("route").Str("generator").Str("generate").Str("client_do")).
			Msg("Error fowarding the request")
		c.JSON(http.StatusBadGateway, gin.H{
			"status":  http.StatusBadGateway,
			"message": "Bad gateway",
		})
		return err
	}
	defer proxyRes.Body.Close()

	resp, err := ioutil.ReadAll(proxyRes.Body)
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Str("host", routeObject.Host).
			Str("request_id", requestID).
			Str("prefix", routeObject.Prefix).
			Str("name", routeObject.Name).
			Str("path", urlPath).
			Str("method", c.Request.Method).
			Str("full_path", c.Request.URL.String()).
			Str("client_ip", c.ClientIP()).
			Array("tags", zerolog.Arr().Str("route").Str("generator").Str("generate").Str("read_all_response")).
			Msg("Error reading the response")
		c.JSON(http.StatusBadGateway, gin.H{
			"status":  http.StatusBadGateway,
			"message": "Bad gateway.",
		})
		return err
	}

	for header, values := range proxyRes.Header {
		for _, value := range values {
			c.Writer.Header().Add(header, value)
		}
	}

	c.Status(proxyRes.StatusCode)
	_, err = c.Writer.Write(resp)
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Str("host", routeObject.Host).
			Str("request_id", requestID).
			Str("prefix", routeObject.Prefix).
			Str("name", routeObject.Name).
			Str("path", urlPath).
			Str("method", c.Request.Method).
			Str("full_path", c.Request.URL.String()).
			Str("client_ip", c.ClientIP()).
			Array("tags", zerolog.Arr().Str("route").Str("generator").Str("generate").Str("writer_write")).
			Msg("Error reading the response")
		return err
	}

	return nil
}

func (g *generator) downStreamPluginMetric(c *gin.Context, routeName, pluginName, path string, startTime time.Time) {
	labels := map[string]string{
		"route_name":        routeName,
		"plugin_name":       pluginName,
		"method":            c.Request.Method,
		"path":              path,
		"status_code":       strconv.Itoa(c.Writer.Status()),
		"status_code_group": strconv.Itoa(((c.Writer.Status() / 100) * 100)),
	}

	g.metric.Observe("routes_downstream_plugin_latency_seconds", float64(time.Since(startTime).Milliseconds()), labels)
}

func (g *generator) downStreamMetric(c *gin.Context, routeName, path string, startTime time.Time) {
	labels := map[string]string{
		"route_name":        routeName,
		"method":            c.Request.Method,
		"path":              path,
		"status_code":       strconv.Itoa(c.Writer.Status()),
		"status_code_group": strconv.Itoa(((c.Writer.Status() / 100) * 100)),
	}

	g.metric.Inc("routes_downstream_hits", labels)
	g.metric.Observe("routes_downstream_latency_seconds", float64(time.Since(startTime).Milliseconds()), labels)
}

func (g *generator) inheritRouterObject(routeObject entity.RouteObject, routePath *entity.RouterPath) {
	if routePath.Auth == "" {
		routePath.Auth = routeObject.Auth
	}
}

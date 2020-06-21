package route

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/journal"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type generator struct {
	routerPath       map[string]entity.RouterPath
	downStreamPlugin []core.DownStreamPlugin
	metric           core.Metric
}

func Generator() core.RouteGenerator {
	return &generator{}
}

func (g *generator) Generate(engine *gin.Engine, metric core.Metric, routeObjects []entity.RouteObject, downStreamPlugin []core.DownStreamPlugin) (errVariable error) {
	g.downStreamPlugin = downStreamPlugin

	g.metric = metric
	g.metric.InjectCounter("routes_downstream_hits", "route_name", "method", "path", "status_code", "status_code_group")
	g.metric.InjectHistogram("routes_downstream_latency_in_ms", "route_name", "method", "path", "status_code", "status_code_group")
	g.metric.InjectHistogram("routes_downstream_plugin_latency_in_ms", "route_name", "plugin_name", "method", "path", "status_code", "status_code_group")

	defer func() {
		if r := recover(); r != nil {
			errVariable = errors.New(fmt.Sprintf("Error generating route because of %v", r))
			journal.Error("Panic error when generating routes", errVariable).
				SetTags("route", "generator", "defer", "panic").
				Log()
		}
	}()

	g.routerPath = map[string]entity.RouterPath{}

	for _, routeObject := range routeObjects {
		for r, routePath := range routeObject.Path {
			g.inheritRouterObject(routeObject, &routePath)

			urlPath := fmt.Sprintf("%s%s", routeObject.Prefix, r)

			g.routerPath[urlPath] = routePath

			journal.Info("Generating routes").
				AddField("host", routeObject.Host).
				AddField("name", routeObject.Name).
				AddField("path", urlPath).
				SetTags("route", "generator", "generate", "url_path").
				Log()

			engine.Any(urlPath, func(c *gin.Context) {
				trackID := uuid.New().String()
				startTime := time.Now()

				g.do(c, urlPath, trackID, routeObject)

				journal.Info("Complete forwarding the request").
					SetTrackId(trackID).
					AddField("host", routeObject.Host).
					AddField("prefix", routeObject.Prefix).
					AddField("name", routeObject.Name).
					AddField("path", urlPath).
					AddField("method", c.Request.Method).
					AddField("full_path", c.Request.URL.String()).
					AddField("client_ip", c.ClientIP()).
					AddField("duration_seconds", time.Since(startTime).Seconds()).
					SetTags("route", "generator", "generate").
					Log()
			})
		}
	}

	return errVariable
}

func (g *generator) do(c *gin.Context, urlPath, trackID string, routeObject entity.RouteObject) {
	proxyReq, err := g.decorateProxyRequest(c, urlPath, trackID, routeObject)
	if err != nil {
		return
	}

	g.decorateHeader(c, proxyReq)

	if err := g.downStreamPluginCallback(c, proxyReq, urlPath, trackID, routeObject); err != nil {
		return
	}

	if err := g.callDownStreamService(c, proxyReq, urlPath, trackID, routeObject); err != nil {
		return
	}

}

func (g *generator) decorateProxyRequest(c *gin.Context, urlPath, trackID string, routeObject entity.RouteObject) (*http.Request, error) {
	var proxyReq *http.Request

	if c.Request.Body != nil {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Malformed request body given by the client",
			})
			journal.Error("Error reading incoming request body", err).
				SetTrackId(trackID).
				AddField("host", routeObject.Host).
				AddField("prefix", routeObject.Prefix).
				AddField("name", routeObject.Name).
				AddField("path", urlPath).
				AddField("method", c.Request.Method).
				AddField("full_path", c.Request.URL.String()).
				AddField("client_ip", c.ClientIP()).
				SetTags("route", "generator", "generate", "read_all_request").
				Log()
			return nil, err
		}

		proxyReq, err = http.NewRequest(c.Request.Method, "", bytes.NewReader(body))
		if err != nil {
			journal.Error("Error creating proxy request", err).
				SetTrackId(trackID).
				AddField("host", routeObject.Host).
				AddField("prefix", routeObject.Prefix).
				AddField("name", routeObject.Name).
				AddField("path", urlPath).
				AddField("method", c.Request.Method).
				AddField("full_path", c.Request.URL.String()).
				AddField("client_ip", c.ClientIP()).
				SetTags("route", "generator", "generate", "new_request").
				Log()
			return nil, err
		}
	} else {
		var err error

		proxyReq, err = http.NewRequest(c.Request.Method, "", nil)
		if err != nil {
			journal.Error("Error creating proxy request", err).
				SetTrackId(trackID).
				AddField("host", routeObject.Host).
				AddField("prefix", routeObject.Prefix).
				AddField("name", routeObject.Name).
				AddField("path", urlPath).
				AddField("method", c.Request.Method).
				AddField("full_path", c.Request.URL.String()).
				AddField("client_ip", c.ClientIP()).
				SetTags("route", "generator", "generate", "new_request").
				Log()
			return nil, err
		}
	}

	proxyReq.URL.Scheme = "http"
	proxyReq.URL.Host = routeObject.Host
	proxyReq.URL.Path = c.Request.URL.Path
	proxyReq.Host = os.Getenv("PROXY_HOST")
	proxyReq.Header.Add("X-Track-ID", trackID)
	proxyReq.Header.Set("X-Real-Ip-Address", c.ClientIP())
	proxyReq.Header.Set("X-Forwarded-For", c.Request.RemoteAddr)

	return proxyReq, nil
}

func (g *generator) decorateHeader(c *gin.Context, proxyReq *http.Request) {
	for header, values := range c.Request.Header {
		for _, value := range values {
			proxyReq.Header.Add(header, value)
		}
	}
}

func (g *generator) downStreamPluginCallback(c *gin.Context, proxyReq *http.Request, urlPath, trackID string, routeObject entity.RouteObject) error {
	for _, plugin := range g.downStreamPlugin {
		startTimePlugin := time.Now()
		if err := plugin.Intervene(c, proxyReq, g.routerPath[urlPath]); err != nil {
			journal.Error("Plugin error", err).
				SetTrackId(trackID).
				AddField("plugin_name", plugin.Name()).
				AddField("plugin_duration_seconds", time.Since(startTimePlugin).Seconds()).
				AddField("host", routeObject.Host).
				AddField("plugin_type", "downstream").
				AddField("prefix", routeObject.Prefix).
				AddField("name", routeObject.Name).
				AddField("path", urlPath).
				AddField("method", c.Request.Method).
				AddField("full_path", c.Request.URL.String()).
				AddField("client_ip", c.ClientIP()).
				SetTags("route", "generator", "generate", "plugin", plugin.Name()).
				Log()
			g.downStreamPluginMetric(c, routeObject.Name, plugin.Name(), urlPath, startTimePlugin)
			return err
		}

		journal.Info("Plugin success").
			SetTrackId(trackID).
			AddField("plugin_name", plugin.Name()).
			AddField("plugin_duration_seconds", time.Since(startTimePlugin).Seconds()).
			SetTags("route", "generator", "generate", "plugin").
			Log()
		g.downStreamPluginMetric(c, routeObject.Name, plugin.Name(), urlPath, startTimePlugin)
	}

	return nil
}

func (g *generator) callDownStreamService(c *gin.Context, proxyReq *http.Request, urlPath, trackID string, routeObject entity.RouteObject) error {
	defer g.downStreamMetric(c, routeObject.Name, urlPath, time.Now())

	client := http.Client{}
	proxyRes, err := client.Do(proxyReq)
	if err != nil {
		journal.Error("Error fowarding the request", err).
			SetTrackId(trackID).
			AddField("host", routeObject.Host).
			AddField("prefix", routeObject.Prefix).
			AddField("name", routeObject.Name).
			AddField("path", urlPath).
			AddField("method", c.Request.Method).
			AddField("full_path", c.Request.URL.String()).
			AddField("client_ip", c.ClientIP()).
			SetTags("route", "generator", "generate", "client_do").
			Log()
		c.JSON(http.StatusBadGateway, gin.H{
			"status":  http.StatusBadGateway,
			"message": "Bad gateway",
		})
		return err
	}
	defer proxyRes.Body.Close()

	resp, err := ioutil.ReadAll(proxyRes.Body)
	if err != nil {
		journal.Error("Error reading the response", err).
			SetTrackId(trackID).
			AddField("host", routeObject.Host).
			AddField("prefix", routeObject.Prefix).
			AddField("name", routeObject.Name).
			AddField("path", urlPath).
			AddField("method", c.Request.Method).
			AddField("full_path", c.Request.URL.String()).
			AddField("client_ip", c.ClientIP()).
			SetTags("route", "generator", "generate", "read_all_response").
			Log()
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
		journal.Error("Error writing the response", err).
			SetTrackId(trackID).
			AddField("host", routeObject.Host).
			AddField("prefix", routeObject.Prefix).
			AddField("name", routeObject.Name).
			AddField("path", urlPath).
			AddField("method", c.Request.Method).
			AddField("full_path", c.Request.URL.String()).
			AddField("client_ip", c.ClientIP()).
			SetTags("route", "generator", "generate", "read_all_response").
			Log()
		journal.Error("Error writing response", err).SetTags("forwader", "core", "writer_write").Log()
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

	g.metric.Observe("routes_downstream_plugin_latency_in_ms", float64(time.Since(startTime).Milliseconds()), labels)
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
	g.metric.Observe("routes_downstream_latency_in_ms", float64(time.Since(startTime).Milliseconds()), labels)
}

func (g *generator) inheritRouterObject(routeObject entity.RouteObject, routePath *entity.RouterPath) {
	if routePath.Auth == "" {
		routePath.Auth = routeObject.Auth
	}
}

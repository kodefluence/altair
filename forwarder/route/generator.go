package route

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/journal"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type generator struct{}

func Generator() core.RouteGenerator {
	return &generator{}
}

func (g *generator) Generate(engine *gin.Engine, routeObjects []entity.RouteObject, downStreamPlugin []core.DownStreamPlugin) (errVariable error) {
	defer func() {
		if r := recover(); r != nil {
			errVariable = errors.New(fmt.Sprintf("Error generating route because of %v", r))
			journal.Error("Panic error when generating routes", errVariable).
				SetTags("route", "generator", "defer", "panic").
				Log()
		}
	}()

	for _, routeObject := range routeObjects {
		for r, routePath := range routeObject.Path {
			urlPath := fmt.Sprintf("%s%s", routeObject.Prefix, r)

			journal.Info("Generating routes").
				AddField("host", routeObject.Host).
				AddField("name", routeObject.Name).
				AddField("path", urlPath).
				SetTags("route", "generator", "generate", "url_path").
				Log()

			engine.Any(urlPath, func(c *gin.Context) {
				var proxyReq *http.Request

				trackID := uuid.New().String()
				startTime := time.Now()

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
							AddField("duration_seconds", time.Since(startTime).Seconds()).
							SetTags("route", "generator", "generate", "read_all_request").
							Log()
						return
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
							AddField("duration_seconds", time.Since(startTime).Seconds()).
							SetTags("route", "generator", "generate", "new_request").
							Log()
						return
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
							AddField("duration_seconds", time.Since(startTime).Seconds()).
							SetTags("route", "generator", "generate", "new_request").
							Log()
						return
					}
				}

				proxyReq.URL.Scheme = "http"
				proxyReq.URL.Host = routeObject.Host
				proxyReq.URL.Path = c.Request.URL.Path
				proxyReq.Host = os.Getenv("PROXY_HOST")
				proxyReq.Header.Set("X-Forwarded-For", c.Request.RemoteAddr)

				for header, values := range c.Request.Header {
					for _, value := range values {
						proxyReq.Header.Add(header, value)
					}
				}

				for _, plugin := range downStreamPlugin {
					startTimePlugin := time.Now()

					if err := plugin.Intervene(c, proxyReq, routePath); err != nil {
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
							AddField("duration_seconds", time.Since(startTime).Seconds()).
							SetTags("route", "generator", "generate", "plugin", plugin.Name()).
							Log()
						return
					}

					journal.Info("Plugin success").
						SetTrackId(trackID).
						AddField("plugin_name", plugin.Name()).
						AddField("plugin_duration_seconds", time.Since(startTimePlugin).Seconds()).
						SetTags("route", "generator", "generate", "plugin").
						Log()
				}

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
						AddField("duration_seconds", time.Since(startTime).Seconds()).
						SetTags("route", "generator", "generate", "client_do").
						Log()
					c.JSON(http.StatusBadGateway, gin.H{
						"status":  http.StatusBadGateway,
						"message": "Bad gateway",
					})
					return
				}

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
						AddField("duration_seconds", time.Since(startTime).Seconds()).
						SetTags("route", "generator", "generate", "read_all_response").
						Log()
					c.JSON(http.StatusBadGateway, gin.H{
						"status":  http.StatusBadGateway,
						"message": "Bad gateway.",
					})
					return
				}

				for header, values := range proxyRes.Header {
					for _, value := range values {
						c.Writer.Header().Add(header, value)
					}
				}

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
						AddField("duration_seconds", time.Since(startTime).Seconds()).
						SetTags("route", "generator", "generate", "read_all_response").
						Log()
					journal.Error("Error writing response", err).SetTags("forwader", "core", "writer_write").Log()
				}

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

func (g *generator) inheritRouterObject() {

}

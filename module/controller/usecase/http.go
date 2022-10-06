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
	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func (ctrl *Controller) InjectHTTP(httpControllers ...module.HttpController) {
	for _, metric := range ctrl.metricController {
		metric.InjectCounter("controller_hits", "method", "path", "status_code", "status_code_group")
		metric.InjectHistogram("controller_elapsed_time_seconds", "method", "path", "status_code", "status_code_group")
	}

	for _, httpController := range httpControllers {
		log.Info().
			Str("path", httpController.Path()).
			Str("path", httpController.Method()).
			Array("tags", zerolog.Arr().Str("altair").Str("controller")).
			Msg("Registering controller")

		ctrl.httpInjector(httpController.Method(), httpController.Path(), func(c *gin.Context) {
			var params string

			requestID := uuid.New()
			startTime := time.Now().UTC()
			c.Set("request_id", requestID)
			c.Set("start_time", startTime)
			ktx := kontext.Fabricate(kontext.WithDefaultContext(c))

			if requestBody := c.Request.Body; requestBody != nil {
				bodyBytes, err := io.ReadAll(requestBody)
				if err == nil {
					c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

					params = string(bodyBytes)
				}
			}

			defer ctrl.httpRecoverFunc(ktx, c, httpController, startTime, params)

			httpController.Control(ktx, c)

			if c.Writer.Status() >= http.StatusBadRequest {
				ctrl.httpLogRequestError(ktx, c, httpController, time.Since(startTime).Milliseconds(), params)
			} else {
				ctrl.httpLogRequestInfo(ktx, c, httpController, time.Since(startTime).Milliseconds(), params)
			}

			ctrl.httpTrackRequest(httpController, time.Since(startTime).Milliseconds(), c.Writer)
		})
	}
}

func (ctrl *Controller) httpRecoverFunc(ktx kontext.Context, c *gin.Context, httpController module.HttpController, startTime time.Time, params string) {
	if err := recover(); err != nil {
		ctrl.httpInternalServerErrorResponse(ktx, c)

		var convertedErr error

		switch err.(type) {
		case error:
			convertedErr = err.(error)
		case string:
			convertedErr = errors.New(err.(string))
		default:
			convertedErr = fmt.Errorf("undefined error: %v", err)
		}

		log.Error().
			Err(convertedErr).
			Stack().
			Interface("request_id", ktx.GetWithoutCheck("request_id")).
			Str("client_ip", c.ClientIP()).
			Stack().
			Array("tags", zerolog.Arr().Str("altair").Str("controller").Str("panic")).
			Msgf("panic received by server")

		ctrl.httpLogRequestError(ktx, c, httpController, time.Since(startTime).Milliseconds(), params)

		ctrl.httpTrackRequest(httpController, time.Since(startTime).Milliseconds(), c.Writer)
		return
	}
}

func (ctrl *Controller) httpInternalServerErrorResponse(ktx kontext.Context, c *gin.Context) {
	response := jsonapi.BuildResponse(ctrl.apiError.InternalServerError(ktx))
	c.JSON(response.HTTPStatus(), response)
}

func (ctrl *Controller) httpLogRequestInfo(ktx kontext.Context, c *gin.Context, httpController module.HttpController, elapsedTime int64, params string) {
	l := log.Info()

	ctrl.httpInstrumentingLog(c, l, params)

	l.Str("relative_path", httpController.Path()).
		Stack().
		Int64("elapsed_time", elapsedTime).
		Interface("request_id", ktx.GetWithoutCheck("request_id")).
		Array("tags", zerolog.Arr().Str("altair").Str("api").Str("ward").Str("endpoint_latency")).
		Msg("Altair endpoint error")

}

func (ctrl *Controller) httpLogRequestError(ktx kontext.Context, c *gin.Context, httpController module.HttpController, elapsedTime int64, params string) {
	l := log.Error()

	ctrl.httpInstrumentingLog(c, l, params)

	l.Str("relative_path", httpController.Path()).
		Err(fmt.Errorf("altair endpoint error: %d", c.Writer.Status())).
		Stack().
		Int64("elapsed_time", elapsedTime).
		Interface("request_id", ktx.GetWithoutCheck("request_id")).
		Array("tags", zerolog.Arr().Str("altair").Str("api").Str("ward").Str("endpoint_latency")).
		Msg("Altair endpoint error")
}

func (ctrl *Controller) httpTrackRequest(httpController module.HttpController, elapsedTime int64, writer gin.ResponseWriter) {
	statusCode := strconv.Itoa(writer.Status())
	statusCodGroup := strconv.Itoa(((writer.Status() / 100) * 100))

	for _, metric := range ctrl.metricController {
		metric.Inc("controller_hits", map[string]string{
			"method":            httpController.Method(),
			"path":              httpController.Path(),
			"status_code":       statusCode,
			"status_code_group": statusCodGroup,
		})

		metric.Observe("controller_elapsed_time_seconds", float64(elapsedTime), map[string]string{
			"method":            httpController.Method(),
			"path":              httpController.Path(),
			"status_code":       statusCode,
			"status_code_group": statusCodGroup,
		})
	}
}

func (ctrl *Controller) httpInstrumentingLog(c *gin.Context, l *zerolog.Event, params string) {
	l.Str("path", c.Request.URL.Path).
		Str("raw_path", c.Request.URL.RawPath).
		Str("raw_query", c.Request.URL.RawQuery).
		Str("user_agent", c.Request.UserAgent()).
		Str("host", c.Request.Host).
		Str("referer", c.Request.Referer()).
		Str("method", c.Request.Method).
		Int("status", c.Writer.Status()).
		Str("client_ip", c.ClientIP()).
		Str("params", params)
}

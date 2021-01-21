package controller

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/codefluence-x/altair/core"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Compile plugin controller
func Compile(engine core.APIEngine, metric core.Metric, ctrl core.Controller) {
	metric.InjectCounter("controller_hits", "method", "path", "status_code", "status_code_group")
	metric.InjectHistogram("controller_elapsed_time_seconds", "method", "path", "status_code", "status_code_group")

	log.Info().
		Str("path", ctrl.Path()).
		Str("path", ctrl.Method()).
		Array("tags", zerolog.Arr().Str("altair").Str("controller")).
		Msg("Registering controller")

	engine.Handle(ctrl.Method(), ctrl.Path(), func(c *gin.Context) {
		var params string

		requestID := uuid.New()
		c.Set("track_id", requestID)
		c.Set("request_id", requestID)
		startTime := time.Now().UTC()

		if requestBody := c.Request.Body; requestBody != nil {
			bodyBytes, err := ioutil.ReadAll(requestBody)
			if err == nil {
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

				params = string(bodyBytes)
			}
		}

		defer recoverFunc(requestID, c, ctrl, metric, startTime, params)

		ctrl.Control(c)

		if c.Writer.Status() >= http.StatusBadRequest {
			logRequestError(requestID, c, ctrl, time.Since(startTime).Milliseconds(), params)
		} else {
			logRequestInfo(requestID, c, ctrl, time.Since(startTime).Milliseconds(), params)
		}

		trackRequest(ctrl, metric, time.Since(startTime).Milliseconds(), c.Writer)
	})
}

func trackRequest(ctrl core.Controller, metric core.Metric, elapsedTime int64, writer gin.ResponseWriter) {
	statusCode := strconv.Itoa(writer.Status())
	statusCodGroup := strconv.Itoa(((writer.Status() / 100) * 100))

	metric.Inc("controller_hits", map[string]string{
		"method":            ctrl.Method(),
		"path":              ctrl.Path(),
		"status_code":       statusCode,
		"status_code_group": statusCodGroup,
	})

	metric.Observe("controller_elapsed_time_seconds", float64(elapsedTime), map[string]string{
		"method":            ctrl.Method(),
		"path":              ctrl.Path(),
		"status_code":       statusCode,
		"status_code_group": statusCodGroup,
	})
}

func recoverFunc(requestID uuid.UUID, c *gin.Context, ctrl core.Controller, metric core.Metric, startTime time.Time, params string) {
	if err := recover(); err != nil {
		internalServerErrorResponse(requestID, c)

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
			Str("request_id", requestID.String()).
			Str("client_ip", c.ClientIP()).
			Stack().
			Array("tags", zerolog.Arr().Str("altair").Str("controller").Str("panic")).
			Msgf("panic received by server")

		logRequestError(requestID, c, ctrl, time.Since(startTime).Milliseconds(), params)
		trackRequest(ctrl, metric, time.Since(startTime).Milliseconds(), c.Writer)
		return
	}
}

func logRequestError(requestID uuid.UUID, c *gin.Context, ctrl core.Controller, elapsedTime int64, params string) {
	l := log.Error()

	instrumentingLog(c, l, params)

	l.Str("relative_path", ctrl.Path()).
		Err(fmt.Errorf("altair endpoint error: %d", c.Writer.Status())).
		Stack().
		Int64("elapsed_time", elapsedTime).
		Str("request_id", requestID.String()).
		Array("tags", zerolog.Arr().Str("altair").Str("api").Str("ward").Str("endpoint_latency")).
		Msg("Altair endpoint error")

}

func instrumentingLog(c *gin.Context, l *zerolog.Event, params string) {
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

func logRequestInfo(requestID uuid.UUID, c *gin.Context, ctrl core.Controller, elapsedTime int64, params string) {
	l := log.Info()

	instrumentingLog(c, l, params)

	l.Str("relative_path", ctrl.Path()).
		Err(fmt.Errorf("altair endpoint error: %d", c.Writer.Status())).
		Stack().
		Int64("elapsed_time", elapsedTime).
		Str("request_id", requestID.String()).
		Array("tags", zerolog.Arr().Str("altair").Str("api").Str("ward").Str("endpoint_latency")).
		Msg("Altair endpoint error")

}

func internalServerErrorResponse(requestID uuid.UUID, c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"message": fmt.Sprintf("Something is not right, help us fix this problem. Contribute to https://github.com/codefluence-x/altair. Or help us by give this code '%s' to the admin of this site.", requestID),
	})
}

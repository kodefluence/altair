package controller

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/journal"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Compile(engine core.APIEngine, metric core.Metric, ctrl core.Controller) {
	metric.InjectCounter("controller_hits", "method", "path", "status_code", "status_code_group")
	metric.InjectHistogram("controller_elapsed_time_seconds", "method", "path", "status_code", "status_code_group")

	journal.Info("Registering controller").
		AddField("path", ctrl.Path()).
		AddField("method", ctrl.Method()).
		SetTags("altair", "controller").
		Log()

	engine.Handle(ctrl.Method(), ctrl.Path(), func(c *gin.Context) {
		var params string

		trackID := uuid.New()
		c.Set("track_id", trackID)
		startTime := time.Now().UTC()

		if requestBody := c.Request.Body; requestBody != nil {
			bodyBytes, err := ioutil.ReadAll(requestBody)
			if err == nil {
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

				params = string(bodyBytes)
			}
		}

		defer recoverFunc(trackID, c, ctrl, metric, startTime, params)

		ctrl.Control(c)

		if c.Writer.Status() >= http.StatusBadRequest {
			logRequestError(trackID, c, ctrl, time.Since(startTime).Milliseconds(), params)
		} else {
			logRequestInfo(trackID, c, ctrl, time.Since(startTime).Milliseconds(), params)
		}

		trackRequest(ctrl, metric, time.Since(startTime).Milliseconds(), c.Writer)
	})
}

func trackRequest(ctrl core.Controller, metric core.Metric, elapsedTime int64, writer gin.ResponseWriter) {
	statusCode := strconv.Itoa(writer.Status())
	statusCodGroup := strconv.Itoa(((writer.Status() / 100) * 100))

	_ = metric.Inc("controller_hits", map[string]string{
		"method":            ctrl.Method(),
		"path":              ctrl.Path(),
		"status_code":       statusCode,
		"status_code_group": statusCodGroup,
	})

	_ = metric.Observe("controller_elapsed_time_seconds", float64(elapsedTime), map[string]string{
		"method":            ctrl.Method(),
		"path":              ctrl.Path(),
		"status_code":       statusCode,
		"status_code_group": statusCodGroup,
	})
}

func recoverFunc(trackID uuid.UUID, c *gin.Context, ctrl core.Controller, metric core.Metric, startTime time.Time, params string) {
	if err := recover(); err != nil {
		internalServerErrorResponse(trackID, c)

		var j journal.Journal

		switch err := err.(type) {
		case error:
			j = journal.Error(fmt.Sprintf("panic received by server. Because of %v", err), err.(error))
		default:
			j = journal.Error(fmt.Sprintf("panic received by server. Because of %v", err), errors.New("altair panic with undefined error"))
		}

		j.AddField("client_ip", c.ClientIP()).
			AddField("traceback", string(debug.Stack())).
			SetTags("altair", "controller", "panic").
			SetTrackId(trackID).
			Log()

		logRequestError(trackID, c, ctrl, time.Since(startTime).Milliseconds(), params)
		trackRequest(ctrl, metric, time.Since(startTime).Milliseconds(), c.Writer)
		return
	}
}

func logRequestError(trackID uuid.UUID, c *gin.Context, ctrl core.Controller, elapsedTime int64, params string) {
	j := journal.Error("Altair endpoint error", fmt.Errorf("altair endpoint error: %d", c.Writer.Status()))

	instrumentingLog(c, j, params)

	j.AddField("relative_path", ctrl.Path()).
		AddField("elapsed_time", elapsedTime).
		SetTags("altair", "api", "ward", "endpoint_latency").
		SetTrackId(trackID).
		Log()
}

func instrumentingLog(c *gin.Context, j journal.Journal, params string) {
	j.AddField("path", c.Request.URL.Path).
		AddField("raw_path", c.Request.URL.RawPath).
		AddField("raw_query", c.Request.URL.RawQuery).
		AddField("user_agent", c.Request.UserAgent()).
		AddField("host", c.Request.Host).
		AddField("referer", c.Request.Referer()).
		AddField("method", c.Request.Method).
		AddField("status", c.Writer.Status()).
		AddField("client_ip", c.ClientIP()).
		AddField("params", params)
}

func logRequestInfo(trackID uuid.UUID, c *gin.Context, ctrl core.Controller, elapsedTime int64, params string) {
	j := journal.Info("Altair endpoint latency and log")

	instrumentingLog(c, j, params)

	j.AddField("relative_path", ctrl.Path()).
		AddField("elapsed_time", elapsedTime).
		SetTags("altair", "api", "ward", "endpoint_latency").
		SetTrackId(trackID).
		Log()

}

func internalServerErrorResponse(trackID uuid.UUID, c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"message": fmt.Sprintf("Something is not right, help us fix this problem. Contribute to https://github.com/codefluence-x/altair. Or help us by give this code '%s' to the admin of this site.", trackID),
	})
}

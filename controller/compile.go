package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/journal"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Compile(engine core.APIEngine, ctrl core.Controller) {

	journal.Info("Registering controller").
		AddField("path", ctrl.Path()).
		AddField("method", ctrl.Method()).
		SetTags("altair", "controller").
		Log()

	engine.Handle(ctrl.Method(), ctrl.Path(), func(c *gin.Context) {
		trackID := uuid.New()
		c.Set("track_id", trackID)
		startTime := time.Now().UTC()

		defer recoverFunc(trackID, c, ctrl, startTime)

		ctrl.Control(c)

		if c.Writer.Status() >= http.StatusBadRequest {
			logRequestError(trackID, c, ctrl, time.Since(startTime).Milliseconds())
		} else {
			logRequestInfo(trackID, c, ctrl, time.Since(startTime).Milliseconds())
		}
	})
}

func recoverFunc(trackID uuid.UUID, c *gin.Context, ctrl core.Controller, startTime time.Time) {
	if err := recover(); err != nil {
		internalServerErrorResponse(trackID, c)

		var j journal.Journal

		switch err.(type) {
		case error:
			j = journal.Error(fmt.Sprintf("panic received by server. Because of %v", err), err.(error))
		case string:
			j = journal.Error(fmt.Sprintf("panic received by server. Because of %v", err), errors.New(err.(string)))
		default:
			j = journal.Error(fmt.Sprintf("panic received by server. Because of %v", err), errors.New("altair panic with undefined error"))
		}

		j.AddField("client_ip", c.ClientIP()).
			SetTags("altair", "controller", "panic").
			SetTrackId(trackID).
			Log()

		logRequestError(trackID, c, ctrl, time.Since(startTime).Milliseconds())
		return
	}
}

func logRequestError(trackID uuid.UUID, c *gin.Context, ctrl core.Controller, elapsedTime int64) {
	var j journal.Journal

	j = journal.Error("Altair endpoint error", errors.New(fmt.Sprintf("altair endpoint error: %d", c.Writer.Status())))

	instrumentingLog(c, j)

	j.AddField("relative_path", ctrl.Path()).
		AddField("elapsed_time", elapsedTime).
		SetTags("altair", "api", "ward", "endpoint_latency").
		SetTrackId(trackID).
		Log()
}

func instrumentingLog(c *gin.Context, j journal.Journal) {
	j.AddField("path", c.Request.URL.Path).
		AddField("raw_path", c.Request.URL.RawPath).
		AddField("raw_query", c.Request.URL.RawQuery).
		AddField("user_agent", c.Request.UserAgent()).
		AddField("host", c.Request.Host).
		AddField("referer", c.Request.Referer()).
		AddField("method", c.Request.Method).
		AddField("status", c.Writer.Status()).
		AddField("client_ip", c.ClientIP())

	if c.Request.Body == nil {
		return
	}

	rawData, err := c.GetRawData()
	if err != nil {
		return
	}

	var payload interface{}
	if err := json.Unmarshal(rawData, &payload); err == nil {
		j.AddField("payload", payload)
	}
}

func logRequestInfo(trackID uuid.UUID, c *gin.Context, ctrl core.Controller, elapsedTime int64) {
	j := journal.Info("Altair endpoint latency and log")

	instrumentingLog(c, j)

	j.AddField("relative_path", ctrl.Path()).
		AddField("elapsed_time", elapsedTime).
		SetTags("altair", "api", "ward", "endpoint_latency").
		SetTrackId(trackID).
		Log()

}

func internalServerErrorResponse(trackID uuid.UUID, c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"message": fmt.Sprintf("Something is not right, help us fix this problem. Contribute to https://github.com/codefluence-x/altair. Or help us by give this code '%s' to the admin of this site.", trackID),
		"meta":    gin.H{"http_status": http.StatusInternalServerError},
	})
}

package model

import (
	"context"
	"time"

	"github.com/codefluence-x/journal"
)

func monitor(ctx context.Context, modelName, query string, f func() error) error {
	startTime := time.Now().UTC()
	err := f()
	elapsedTime := time.Since(startTime).Milliseconds()
	if err != nil {
		journal.Error("Failed executing query", err).
			SetTrackId(ctx.Value("track_id")).
			AddField("duration_in_ms", elapsedTime).
			AddField("model_name", modelName).
			AddField("query", query).
			SetTags("model", "monitor").
			Log()
		return err
	}

	journal.Info("Complete executing query").
		SetTrackId(ctx.Value("track_id")).
		AddField("duration_in_ms", elapsedTime).
		AddField("model_name", modelName).
		AddField("query", query).
		SetTags("model", "monitor").
		Log()
	return nil
}

package model

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func monitor(ctx context.Context, modelName, query string, f func() error) error {
	startTime := time.Now().UTC()
	err := f()
	elapsedTime := time.Since(startTime).Milliseconds()
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Interface("request_id", ctx.Value("request_id")).
			Int64("duration_in_ms", elapsedTime).
			Str("model_name", modelName).
			Str("query", query).
			Array("tags", zerolog.Arr().Str("model").Str("monitor")).
			Msg("Failed executing query")
		return err
	}

	log.Info().
		Interface("request_id", ctx.Value("request_id")).
		Int64("duration_in_ms", elapsedTime).
		Str("model_name", modelName).
		Str("query", query).
		Array("tags", zerolog.Arr().Str("model").Str("monitor")).
		Msg("Complete executing query")
	return nil
}

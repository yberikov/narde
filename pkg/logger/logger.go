package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
	"time"
)

func ParseEnvLoggerEnv(value string) zerolog.Level {
	v := strings.ToLower(value)

	var level zerolog.Level

	switch v {
	case "info":
		level = zerolog.InfoLevel
	case "warn":
		level = zerolog.WarnLevel
	case "error":
		level = zerolog.ErrorLevel
	case "debug":
		level = zerolog.DebugLevel
	case "trace":
		level = zerolog.TraceLevel
	default:
		log.Warn().Msgf("Unknown logging level: %s", value)

		level = zerolog.InfoLevel
	}

	return level
}

func InitRootLogger(logForcePlainText bool, globalLevel zerolog.Level, serviceName string) zerolog.Logger {
	if !logForcePlainText {
		log.Logger = log.Output(
			zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: time.RFC3339,
			},
		)
	}

	zerolog.SetGlobalLevel(globalLevel)

	log.Logger = log.With().
		Str("service", serviceName).
		Caller().
		Logger()
	zerolog.DefaultContextLogger = &log.Logger

	return log.Logger
}

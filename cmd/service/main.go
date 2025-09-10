package main

import (
	"github.com/caarlos0/env/v11"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"narde/internal/config"
	"narde/pkg/logger"
)

func main() {

	cfg, err := env.ParseAs[config.Config]()
	if err != nil {
		log.Fatal().Err(err).Msg("Can't parse config")
	}

	rootLogger := logger.InitRootLogger(cfg.Main.LogForcePlainText, logger.ParseEnvLoggerEnv(cfg.Main.LogLevel), "backend")
	rootLogger.Info().Msg("Test")
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.Run() // listen and serve on 0.0.0.0:8080
}

package main

import (
	"github.com/caarlos0/env/v11"
	"github.com/rs/zerolog/log"
	"narde/internal/config"
	"narde/internal/service"
	"narde/internal/transport/http"
	"narde/internal/transport/http/handlers"
	"narde/internal/transport/http/router"
	"narde/pkg/jwt"
	"narde/pkg/logger"
)

func main() {

	cfg, err := env.ParseAs[config.Config]()
	if err != nil {
		log.Fatal().Err(err).Msg("Can't parse config")
	}

	rootLogger := logger.InitRootLogger(cfg.Main.LogForcePlainText, logger.ParseEnvLoggerEnv(cfg.Main.LogLevel), "backend")

	//auth := middleware.NewAuth(jwt.MustParser())
	authHandler := handlers.NewAuthHandler(service.NewAuthService(jwt.MustGenerator()))

	server := http.NewServer(":8080", router.NewRouter(authHandler.Router()))
	rootLogger.Info().Msg("Server is running")
	server.Run()

}

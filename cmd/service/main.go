package main

import (
	"context"
	"github.com/caarlos0/env/v11"
	"github.com/rs/zerolog/log"
	"narde/internal/config"
	"narde/internal/repository"
	"narde/internal/service"
	"narde/internal/transport/http"
	"narde/internal/transport/http/handlers"
	"narde/internal/transport/http/router"
	websocket2 "narde/internal/transport/websocket"
	websocket "narde/internal/transport/websocket/handlers"
	"narde/pkg/jwt"
	"narde/pkg/logger"
)

func main() {

	cfg, err := env.ParseAs[config.Config]()
	if err != nil {
		log.Fatal().Err(err).Msg("Can't parse config")
	}

	rootLogger := logger.InitRootLogger(cfg.Main.LogForcePlainText, logger.ParseEnvLoggerEnv(cfg.Main.LogLevel), "backend")
	ctx := context.Background()

	db, err := setupPostgresConnection(ctx, cfg.Postgresdb.Main)
	if err != nil {
		log.Fatal().Err(err).Msgf("Can't setup postgres connection to %s database", cfg.Postgresdb.Main.DBName)
	}
	defer db.Close()

	if err = runMigrations(cfg.Postgresdb.Main, db); err != nil {
		log.Fatal().Err(err).Msgf("Failed to apply migrations to %s PostgreSQL database", cfg.Postgresdb.Main.DBName)
	}

	//auth := middleware.NewAuth(jwt.MustParser())
	authRepository := repository.NewAuthRepository(db)
	authHandler := handlers.NewAuthHandler(service.NewAuthService(authRepository, jwt.MustGenerator()))
	hub := websocket2.NewHub()
	go hub.Run()
	hubHandler := websocket.NewHubHander(jwt.MustParser(), hub)
	server := http.NewServer(":8080",
		router.NewRouter(authHandler.Router()),
		router.NewRouter(hubHandler.Router()))
	rootLogger.Info().Msg("Server is running")
	server.Run()

}

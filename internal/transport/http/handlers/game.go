package handlers

import (
	"github.com/gofiber/fiber/v3"
	"narde/internal/transport/http/router"
)

type (
	GameHandler struct {
		makeRouter router.MakeRouter
	}
)

func NewGameHandler() *GameHandler {
	gameHandler := &GameHandler{}
	makeRouter := func(router fiber.Router) {
		g := router.Group("/games")
		g.Get("/ping", gameHandler.Pong)
	}

	gameHandler.makeRouter = makeRouter

	return gameHandler
}

func (h *GameHandler) Router() router.MakeRouter {
	return h.makeRouter
}

func (*GameHandler) Pong(ctx fiber.Ctx) error {
	ctx.Write([]byte("Pong"))
	return nil
}

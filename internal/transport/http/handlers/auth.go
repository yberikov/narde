package handlers

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
	"narde/internal/transport/http/handlers/model"
	"narde/internal/transport/http/router"
	"net/http"
)

type (
	AuthHandler struct {
		makeRouter  router.MakeRouter
		authService AuthService
	}

	AuthService interface {
		Login(ctx context.Context) (*model.TokenResponse, error)
	}
)

func NewAuthHandler(authService AuthService) *AuthHandler {
	authHandler := &AuthHandler{authService: authService}
	makeRouter := func(router fiber.Router) {
		g := router.Group("/auth")
		g.Post("/register", authHandler.Register)
		g.Post("/login", authHandler.Login)
	}

	authHandler.makeRouter = makeRouter

	return authHandler
}

func (h *AuthHandler) Router() router.MakeRouter {
	return h.makeRouter
}

func (h *AuthHandler) Register(ctx fiber.Ctx) error {
	ctx.Write([]byte("Pong"))

	return nil
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	logger := zerolog.Ctx(c.RequestCtx())
	var request model.LoginRequest
	if err := c.Bind().Body(&request); err != nil {
		logger.Error().Err(err).Send()
		return c.Status(http.StatusBadRequest).JSON(model.NewErrorResponse("Cannot parse body"))
	}

	response, err := h.authService.Login(c.RequestCtx())
	switch {
	case err != nil:
		return c.Status(http.StatusInternalServerError).JSON(model.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(response)
}

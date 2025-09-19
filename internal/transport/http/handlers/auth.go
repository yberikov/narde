package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
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
		Register(ctx context.Context, req *model.RegisterRequest) error
		Login(ctx context.Context, req *model.LoginRequest) (*model.TokenResponse, error)
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

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	logger := zerolog.Ctx(c.Context())
	var request model.RegisterRequest
	if err := c.BodyParser(&request); err != nil {
		logger.Error().Err(err).Send()
		return c.Status(http.StatusBadRequest).JSON(model.NewErrorResponse("Cannot parse body"))
	}

	err := h.authService.Register(c.Context(), &request)
	switch {
	case err != nil:
		return c.Status(http.StatusInternalServerError).JSON(model.ErrorResponse{Error: err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	logger := zerolog.Ctx(c.Context())
	var request model.LoginRequest
	if err := c.BodyParser(&request); err != nil {
		logger.Error().Err(err).Send()
		return c.Status(http.StatusBadRequest).JSON(model.NewErrorResponse("Cannot parse body"))
	}

	response, err := h.authService.Login(c.Context(), &request)
	switch {
	case err != nil:
		return c.Status(http.StatusInternalServerError).JSON(model.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(response)
}

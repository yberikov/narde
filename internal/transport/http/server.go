package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

type (
	Router interface {
		Init(app *fiber.App)
	}
	Server struct {
		fiber   *fiber.App
		address string
		routers []Router
	}
)

func NewServer(address string, routers ...Router) *Server {
	s := &Server{
		fiber:   fiber.New(),
		address: address,
		routers: routers,
	}

	s.fiber.Use(
		cors.New(cors.Config{
			AllowOrigins: []string{"*"},
		}),
	)

	for idx := range s.routers {
		s.routers[idx].Init(s.fiber)
	}

	return s
}

func (s *Server) Run() error {
	return s.fiber.Listen(s.address)
}

func (s *Server) ShutDown() error {
	return s.fiber.Shutdown()
}

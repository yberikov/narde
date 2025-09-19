package router

import "github.com/gofiber/fiber/v2"

type (
	MakeRouter func(router fiber.Router)

	Router struct {
		makeRouter MakeRouter
	}
)

func NewRouter(makeFunc MakeRouter) *Router {
	return &Router{
		makeRouter: makeFunc,
	}
}

func (r *Router) Init(app *fiber.App) {
	mainGroup := app.Group("api/v1")

	r.makeRouter(mainGroup)
}

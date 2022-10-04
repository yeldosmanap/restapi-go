package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	"gorest-api/docs"
	"gorest-api/internal/service"

	_ "github.com/swaggo/files"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutesFiber(app *fiber.App) *fiber.App {
	docs.SwaggerInfo.BasePath = "/"

	auth := app.Group("/auth")
	{
		auth.Post("sign-up", h.signUp)
		auth.Post("sign-in", h.signIn)
	}

	rest := app.Group("/api/v1")
	{
		projects := rest.Group("/projects")
		{
			projects.Post("/", h.createProject)
			projects.Get("/", h.getAllProjects)
			projects.Get("/:title", h.getProjectByTitle)
			projects.Put("/:id", h.updateProject)
			projects.Delete("/:id", h.deleteProject)
		}
	}

	app.Get("/swagger/*", swagger.HandlerDefault)

	return app
}

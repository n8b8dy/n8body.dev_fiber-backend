package routes

import (
	"github.com/gofiber/fiber/v2"
	"n8body.dev/fiber-backend/app/controllers"
)

// PublicRoutes sets up the public routes
func PublicRoutes(app *fiber.App, controller *controllers.MainController) {
	api := app.Group("/api")

	api.Get("/healthcheck", controller.Healthcheck)

	api.Get("/owner", controller.GetOwner)
	api.Get("/content/home", controller.GetHomeSections)

	api.Post("/login", controller.Login)

	api.Post("/message", controller.CreateMessage)
}

package routes

import (
	"github.com/gofiber/fiber/v2"
	"n8body.dev/fiber-backend/app/controllers"
	"n8body.dev/fiber-backend/pkg/middleware"
)

// PrivateRoutes sets up the public routes
func PrivateRoutes(app *fiber.App, controller *controllers.MainController) {
	api := app.Group("/api", middleware.JWTProtected())

	api.Get("/admin", controller.GetAmIAdmin)

	api.Get("/user/me", controller.GetUserMe)

	api.Get("/messages", controller.GetMessages)
	api.Get("/message/:id", controller.GetMessage)
	api.Delete("/message/:id", controller.DeleteMessage)
}

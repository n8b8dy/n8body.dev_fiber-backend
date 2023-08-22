package main

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
	"n8body.dev/fiber-backend/app/controllers"
	"n8body.dev/fiber-backend/app/dtos"
	"n8body.dev/fiber-backend/pkg/middleware"
	"n8body.dev/fiber-backend/platform/database"

	"n8body.dev/fiber-backend/pkg/routes"
	"n8body.dev/fiber-backend/pkg/utils"
)

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			return ctx.Status(code).JSON(fiber.Map{
				"error": dtos.ErrorResponseDTO{
					Status: code,
					Reason: err.Error(),
				},
			})
		}},
	)

	db, err := database.OpenDBConnection()
	if err != nil {
		log.Fatal(err)
	}

	controller := controllers.NewMainController(db)

	middleware.FiberMiddleware(app)

	routes.PublicRoutes(app, controller)
	routes.PrivateRoutes(app, controller)

	if err := utils.StartServer(app); err != nil {
		log.Fatalln(err)
	}
}

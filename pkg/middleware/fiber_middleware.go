package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// FiberMiddleware provides built-in Fiber's middlewares.
func FiberMiddleware(app *fiber.App) {
	app.Use(
		cors.New(cors.Config{
			AllowOrigins: os.Getenv("ALLOWED_ORIGINS"),
		}),
		helmet.New(),
		logger.New(),
		recover.New(),
	)
}

package utils

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
)

func StartServer(app *fiber.App) error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	address := fmt.Sprintf(":%s", port)

	return app.Listen(address)
}

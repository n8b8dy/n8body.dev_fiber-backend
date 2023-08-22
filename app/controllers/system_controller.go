package controllers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SystemController struct {
	*gorm.DB
}

func (controller *SystemController) Healthcheck(ctx *fiber.Ctx) error {
	db, err := controller.DB.DB()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = db.Ping()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": nil,
	})
}

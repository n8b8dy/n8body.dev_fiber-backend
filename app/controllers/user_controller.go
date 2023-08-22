package controllers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"n8body.dev/fiber-backend/app/dtos"
	"n8body.dev/fiber-backend/app/models"
	"n8body.dev/fiber-backend/pkg/utils"
)

type UserController struct {
	*gorm.DB
}

// GetUserMe accepts User ID via JWT and responds with { error: string | nil, user: UserResponseDTO | nil }
func (controller *UserController) GetUserMe(ctx *fiber.Ctx) error {
	userID, err := utils.GetUserIDFromContextWithJWT(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusInternalServerError,
				Reason: err.Error(),
			},
			"user": nil,
		})
	}

	var user = models.User{}

	if err := controller.Model(&models.User{}).Find(&user, userID).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusNotFound,
				Reason: err.Error(),
			},
			"user": nil,
		})
	} else if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusInternalServerError,
				Reason: err.Error(),
			},
			"user": nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": nil,
		"user": dtos.UserResponseDTO{
			BaseResponseDTO: dtos.BaseResponseDTO(user.BaseModel),
			Username:        user.Username,
			Email:           user.Email,
			Admin:           user.Admin,
		},
	})
}

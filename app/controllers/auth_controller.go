package controllers

import (
	"errors"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	jtoken "github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"n8body.dev/fiber-backend/app/dtos"
	"n8body.dev/fiber-backend/app/models"
	"n8body.dev/fiber-backend/pkg/utils"
)

type AuthController struct {
	*gorm.DB
}

// Login accepts LoginRequestDTO and responds with { error: string | nil, token: string | nil }
func (controller *AuthController) Login(ctx *fiber.Ctx) error {
	var requestBody dtos.LoginRequestDTO

	if err := ctx.BodyParser(&requestBody); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusBadRequest,
				Reason: err.Error(),
			},
			"token": nil,
		})
	}

	validate := utils.NewValidator()

	if err := validate.Struct(requestBody); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusBadRequest,
				Reason: err.Error(),
			},
			"token": nil,
		})
	}

	var user models.User

	if err := controller.Where("email = ?", requestBody.Email).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusNotFound,
				Reason: err.Error(),
			},
			"token": nil,
		})
	} else if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusInternalServerError,
				Reason: err.Error(),
			},
			"token": nil,
		})
	}

	if !utils.CompareHashAndPassword(user.Password, requestBody.Password) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusUnauthorized,
				Reason: "user with the given credentials not found",
			},
			"token": nil,
		})
	}

	claims := jtoken.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"exp":      time.Now().Add(time.Hour * 24 * 1).Unix(),
	}

	tokenWithClaims := jtoken.NewWithClaims(jtoken.SigningMethodHS512, claims)

	token, err := tokenWithClaims.SignedString([]byte(os.Getenv("JWT_PRIVATE_KEY")))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusInternalServerError,
				Reason: err.Error(),
			},
			"token": nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": nil,
		"token": token,
	})
}

// GetAmIAdmin accepts User ID via JWT and responds with { error: string | nil, admin: bool | nil }
func (controller *AuthController) GetAmIAdmin(ctx *fiber.Ctx) error {
	userID, err := utils.GetUserIDFromContextWithJWT(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusInternalServerError,
				Reason: err.Error(),
			},
			"admin": nil,
		})
	}

	var admin bool

	if err := controller.Model(&models.User{}).Select("admin").First(&admin, userID).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusNotFound,
				Reason: err.Error(),
			},
			"admin": nil,
		})
	} else if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusInternalServerError,
				Reason: err.Error(),
			},
			"admin": nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": nil,
		"admin": admin,
	})
}

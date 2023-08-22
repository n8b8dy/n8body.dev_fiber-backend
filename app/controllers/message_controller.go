package controllers

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"n8body.dev/fiber-backend/app/dtos"
	"n8body.dev/fiber-backend/app/models"
	"n8body.dev/fiber-backend/pkg/utils"
)

type MessageController struct {
	*gorm.DB
}

// GetMessages receives sort_field, sort_order via queries and responds with { error: string | nil, count: number | nil, messages: Array<MessageResponseDTO> | nil }
func (controller *MessageController) GetMessages(ctx *fiber.Ctx) error {
	userID, err := utils.GetUserIDFromContextWithJWT(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusInternalServerError,
				Reason: err.Error(),
			},
			"count":    nil,
			"messages": nil,
		})
	}

	queries := ctx.Queries()
	sortField := queries["sort_field"]
	sortOrder := queries["sort_order"]

	if sortField == "" {
		sortField = "created_at"
	}
	if sortOrder == "" {
		sortOrder = "desc"
	}

	var admin bool

	if err := controller.Model(&models.User{}).Select("admin").First(&admin, userID).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusUnauthorized,
				Reason: err.Error(),
			},
			"count":    nil,
			"messages": nil,
		})
	} else if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusInternalServerError,
				Reason: err.Error(),
			},
			"count":    nil,
			"messages": nil,
		})
	}

	if !admin {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusForbidden,
				Reason: "you don't have access to this information",
			},
			"count":    nil,
			"messages": nil,
		})
	}

	var messages []models.Message

	if err := controller.Order(fmt.Sprintf("%s %s", sortField, sortOrder)).Find(&messages).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusNotFound,
				Reason: err.Error(),
			},
			"count":    0,
			"messages": fiber.Map{},
		})
	} else if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusInternalServerError,
				Reason: err.Error(),
			},
			"count":    nil,
			"messages": nil,
		})
	}

	messagesResponseDTOs := make([]dtos.MessageResponseDTO, len(messages))

	for i, v := range messages {
		messagesResponseDTOs[i] = dtos.MessageResponseDTO{
			BaseResponseDTO: dtos.BaseResponseDTO(v.BaseModel),
			Text:            v.Text,
			Username:        v.Username,
			Email:           v.Email,
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":    nil,
		"count":    len(messagesResponseDTOs),
		"messages": messagesResponseDTOs,
	})
}

// GetMessage accepts ID via params responds with { error: string | nil, message: MessageResponseDTO | nil }
func (controller *MessageController) GetMessage(ctx *fiber.Ctx) error {
	userID, err := utils.GetUserIDFromContextWithJWT(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusInternalServerError,
				Reason: err.Error(),
			},
			"message": nil,
		})
	}

	var admin bool

	if err := controller.Model(&models.User{}).Select("admin").First(&admin, userID).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusUnauthorized,
				Reason: err.Error(),
			},
			"message": nil,
		})
	} else if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusInternalServerError,
				Reason: err.Error(),
			},
			"message": nil,
		})
	}

	if !admin {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusForbidden,
				Reason: "you don't have access to this information",
			},
			"message": nil,
		})
	}

	messageID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusBadRequest,
				Reason: err.Error(),
			},
			"message": nil,
		})
	}

	var message models.Message

	if err := controller.First(&message, messageID).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusNotFound,
				Reason: err.Error(),
			},
			"message": nil,
		})
	} else if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusInternalServerError,
				Reason: err.Error(),
			},
			"message": nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": nil,
		"message": dtos.MessageResponseDTO{
			BaseResponseDTO: dtos.BaseResponseDTO(message.BaseModel),
			Text:            message.Text,
			Username:        message.Username,
			Email:           message.Email,
		},
	})
}

// CreateMessage accepts MessageRequestDTO as body and responds with { error: string | nil, message: MessageResponseDTO | nil }
func (controller *MessageController) CreateMessage(ctx *fiber.Ctx) error {
	var requestBody dtos.MessageRequestDTO

	if err := ctx.BodyParser(&requestBody); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusBadRequest,
				Reason: err.Error(),
			},
			"message": nil,
		})
	}

	validate := utils.NewValidator()

	if err := validate.Struct(&requestBody); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusBadRequest,
				Reason: err.Error(),
			},
			"message": nil,
		})
	}

	message := models.Message{
		BaseModel: models.BaseModel{},
		Text:      requestBody.Text,
		Username:  requestBody.Username,
		Email:     requestBody.Email,
	}

	if err := controller.Create(&message).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusInternalServerError,
				Reason: err.Error(),
			},
			"message": nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": nil,
		"message": dtos.MessageResponseDTO{
			BaseResponseDTO: dtos.BaseResponseDTO(message.BaseModel),
			Text:            message.Text,
			Username:        message.Username,
			Email:           message.Email,
		},
	})
}

// DeleteMessage accepts ID via params responds with { error: string | nil, message: Message | nil }
func (controller *MessageController) DeleteMessage(ctx *fiber.Ctx) error {
	userID, err := utils.GetUserIDFromContextWithJWT(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusInternalServerError,
				Reason: err.Error(),
			},
			"message": nil,
		})
	}

	var admin bool

	if err := controller.Model(&models.User{}).Select("admin").First(&admin, userID).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusUnauthorized,
				Reason: err.Error(),
			},
			"message": nil,
		})
	} else if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusInternalServerError,
				Reason: err.Error(),
			},
			"message": nil,
		})
	}

	if !admin {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusForbidden,
				Reason: "you don't have access to this information",
			},
			"message": nil,
		})
	}

	messageID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusBadRequest,
				Reason: err.Error(),
			},
			"message": nil,
		})
	}

	var message models.Message

	if err := controller.Clauses(clause.Returning{}).Delete(&message, messageID).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusNotFound,
				Reason: err.Error(),
			},
			"message": nil,
		})
	} else if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusInternalServerError,
				Reason: err.Error(),
			},
			"message": nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": nil,
		"message": dtos.MessageResponseDTO{
			BaseResponseDTO: dtos.BaseResponseDTO(message.BaseModel),
			Text:            message.Text,
			Username:        message.Username,
			Email:           message.Email,
		},
	})
}

package controllers

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"n8body.dev/fiber-backend/app/dtos"
	"n8body.dev/fiber-backend/app/models"
)

type ContentController struct {
	*gorm.DB
}

// GetOwner responds with { error: nil, owner: Owner }
func (*ContentController) GetOwner(ctx *fiber.Ctx) error {
	owner := dtos.OwnerResponseDTO{
		UserResponseDTO: dtos.UserResponseDTO{
			BaseResponseDTO: dtos.BaseResponseDTO{
				ID:        uuid.UUID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 11, 3, 176, 59, 24, 120},
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
				DeletedAt: gorm.DeletedAt{},
			},
			Username: "n8body",
			Email:    "contact@n8body.dev",
			Admin:    true,
		},
		Status: "searching for a job",
		Age:    17,
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": nil,
		"owner": owner,
	})
}

// GetHomeSections responds with { error: string | nil, count: number | nil, sections: Array<SectionResponseDTO> }
func (controller *ContentController) GetHomeSections(ctx *fiber.Ctx) error {
	var sections []models.Section

	if err := controller.Order("position").Find(&sections).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusNotFound,
				Reason: err.Error(),
			},
			"count":    nil,
			"sections": nil,
		})
	} else if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": dtos.ErrorResponseDTO{
				Status: fiber.StatusInternalServerError,
				Reason: err.Error(),
			},
			"count":    nil,
			"sections": nil,
		})
	}

	sectionResponseDTOs := make([]dtos.SectionResponseDTO, len(sections))

	for i, v := range sections {
		sectionResponseDTOs[i] = dtos.SectionResponseDTO{
			BaseResponseDTO: dtos.BaseResponseDTO(v.BaseModel),
			Position:        v.Position,
			Heading:         v.Heading,
			Paragraphs:      v.Paragraphs,
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":    nil,
		"count":    len(sectionResponseDTOs),
		"sections": sectionResponseDTOs,
	})
}

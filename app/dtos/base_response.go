package dtos

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseResponseDTO struct {
	ID        uuid.UUID      `json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
}

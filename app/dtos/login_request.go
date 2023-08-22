package dtos

type LoginRequestDTO struct {
	Email    string `json:"email" validate:"required,email,min=3,max=255"`
	Password string `json:"password" validate:"required"`
}

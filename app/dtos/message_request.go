package dtos

type MessageRequestDTO struct {
	Text     string `json:"text" validate:"required,min=3,max=2047"`
	Username string `json:"username" validate:"required,min=3,max=255"`
	Email    string `json:"email" validate:"required,email,min=3,max=255"`
}

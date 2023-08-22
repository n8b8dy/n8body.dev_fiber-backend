package dtos

type UserRequestDTO struct {
	Username string `json:"username" validate:"required,min=3,max=255"`
	Email    string `json:"email" validate:"required,email,min=3,max=255"`
	Password string `json:"password" validate:"required"`
	Admin    bool   `json:"admin" validate:"required,boolean"`
}

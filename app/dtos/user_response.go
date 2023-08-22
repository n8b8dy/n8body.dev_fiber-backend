package dtos

type UserResponseDTO struct {
	BaseResponseDTO
	Username string `json:"username"`
	Email    string `json:"email"`
	Admin    bool   `json:"admin"`
}

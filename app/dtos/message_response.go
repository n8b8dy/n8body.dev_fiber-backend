package dtos

type MessageResponseDTO struct {
	BaseResponseDTO
	Text     string `json:"text"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

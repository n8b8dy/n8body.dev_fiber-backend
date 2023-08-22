package dtos

type ErrorResponseDTO struct {
	Status int    `json:"status"`
	Reason string `json:"reason"`
}

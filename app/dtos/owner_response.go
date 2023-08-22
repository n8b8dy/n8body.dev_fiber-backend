package dtos

type OwnerResponseDTO struct {
	UserResponseDTO
	Status string `json:"status"`
	Age    uint8  `json:"age"`
}

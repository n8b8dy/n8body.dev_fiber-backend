package dtos

type SectionResponseDTO struct {
	BaseResponseDTO
	Position   int8     `json:"position"`
	Heading    string   `json:"heading"`
	Paragraphs []string `json:"paragraphs"`
}

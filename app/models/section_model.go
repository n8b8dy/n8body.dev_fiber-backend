package models

import "github.com/lib/pq"

type Section struct {
	BaseModel
	Position   int8
	Heading    string
	Paragraphs pq.StringArray `gorm:"type:text"`
}

package models

type Message struct {
	BaseModel
	Text     string
	Username string
	Email    string
}

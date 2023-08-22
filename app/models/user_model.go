package models

type User struct {
	BaseModel
	Username string
	Email    string
	Password string
	Admin    bool
}

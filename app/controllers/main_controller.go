package controllers

import "gorm.io/gorm"

type MainController struct {
	*SystemController
	*AuthController
	*MessageController
	*ContentController
	*UserController
}

func NewMainController(db *gorm.DB) *MainController {
	return &MainController{
		SystemController:  &SystemController{db},
		AuthController:    &AuthController{db},
		MessageController: &MessageController{db},
		ContentController: &ContentController{db},
		UserController:    &UserController{db},
	}
}

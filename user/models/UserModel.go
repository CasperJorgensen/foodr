package user_models

import (
	misc "foodr/misc/models"
)

type UserDTO struct {
	misc.Base
	Fullname string `json:"fullname"`
	Username string `json:"username" binding:"required"`
	Email    string `gorm: "not null; unique" json"email" binding:"required"`
	Password string `json:"password"`
}

type UserLoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

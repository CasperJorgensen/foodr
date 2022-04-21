package user

import (
	"github.com/gin-gonic/gin"
	"foodr/misc/models"
)

func FullnameCannotBeEmpty() *gin.H {
	return &gin.H{"error": misc.NewError(
		"Fullname cannot be null or empty",
		misc.FullnameCannotBeNullOrEmpty,
		true,
	)}
}

func UsernameCannotBeEmpty() *gin.H {
	return &gin.H{"error": misc.NewError(
		"User cannot be null or empty",
		misc.UsernameCannotBeNullOrEmpty,
		true,
	)}
}

func UsernameAlreadyExist() *gin.H {
	return &gin.H{"error": misc.NewError(
		"Username already exist",
		misc.UsernameAlreadyExist,
		true,
	)}
}

func EmailAlreadyExist() *gin.H {
	return &gin.H{"error": misc.NewError(
		"Email already exist",
		misc.EmailAlreadyExist,
		true,
	)}
}
package recipe

import (
	"github.com/gin-gonic/gin"
	"foodr/misc/models"
)

func RecipeNotFoundError() *gin.H {
	return &gin.H{"error": misc.NewError(
		"Recipe not found",
		misc.RecipeNotFound,
		true,
	)}
}
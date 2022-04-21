package recipe_model

import (
	misc "foodr/misc/models"
	user_models "foodr/user/models"

	"github.com/google/uuid"
)

type RecipeDTO struct {
	misc.Base
	Title   string              `gorm: "not null; unique" json:"title"`
	OwnerID uuid.UUID           `gorm: "foreign_key" json:"owner_id"`
	Owner   user_models.UserDTO `json:"owner"`
}

type CreateRecipeDTO struct {
	misc.Base
	Title string `gorm: "not null; unique" json:"title"`
}

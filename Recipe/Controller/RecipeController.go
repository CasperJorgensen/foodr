package recipe

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/morkid/paginate"

	token_auth "foodr/Authentication/Controller"
	recipe "foodr/Recipe"
	recipe_model "foodr/Recipe/Models"
	"foodr/configuration"
	user_model "foodr/user/models"
)

func CreateRecipe(c *gin.Context) {
	var recipe recipe_model.CreateRecipeDTO
	var user user_model.UserDTO

	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "invalid json")
		return
	}

	metadata, err := token_auth.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	userid, err := token_auth.FetchAuth(metadata)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	if err := configuration.DB.Where("id = ?", userid).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "User not found"})
		return
	}

	useruuid, err := uuid.Parse(userid)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Something went wrong"})
	}

	recipy := recipe_model.RecipeDTO{Title: recipe.Title, OwnerID: useruuid}

	configuration.DB.Create(&recipy)
	c.JSON(http.StatusCreated, gin.H{"item": recipy})
}

func GetRecipe(c *gin.Context) {
	var recipeDTO recipe_model.RecipeDTO
	// if err := configuration.DB.Where("id = ?", c.Param("id")).First(&recipe).Error; err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"Error": "Recipe not found"})
	// 	return
	// }

	// if err := configuration.DB.Joins("JOIN user_dtos on user_dtos.id = recipe_dtos.owner_id").Find(&recipeResponse).Error; err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"Error": "Something went wrong"})
	// 	fmt.Println(recipeResponse)
	// 	return
	// }

	// if err := configuration.DB.Where("id = ?", recipe.OwnerID).First(&user).Error; err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"Error": "User not found"})
	// 	return
	// }

	if err := configuration.DB.
		Where("recipe_dtos.id = ?", c.Param("id")).
		Joins("JOIN user_dtos on recipe_dtos.owner_id=user_dtos.id::text").
		Preload("Owner").
		Find(&recipeDTO).Error; err != nil {
		c.JSON(
			http.StatusNotFound,
			recipe.RecipeNotFoundError())
		return
	}

	// configuration.DB.First(&recipeDTO, "id = ?", c.Param("id"))

	c.JSON(http.StatusOK, gin.H{"item": recipeDTO})
}

func GetAllRecipeFromUser(c *gin.Context) {
	var recipeDTO []*recipe_model.RecipeDTO
	pg := paginate.New()

	// model := db.
	// 	Where("recipe_dtos.owner_id = ?", c.Param("id")).
	// 	Preload("Owner").
	// 	Find(&recipeDTO)

	// model := configuration.DB.Where("recipe_dtos.owner_id = ?", c.Param("id"))
	model := configuration.DB.
		Where("recipe_dtos.owner_id = ?", c.Param("id")).
		Joins("JOIN user_dtos on recipe_dtos.owner_id=user_dtos.id::text").
		Preload("Owner").
		Find(&recipeDTO)

	c.JSON(http.StatusOK, pg.With(model).Request(c.Request).Response(&recipeDTO))
}

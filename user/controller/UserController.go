package user_controller

import (
	"fmt"
	"net/http"

	token_auth "foodr/Authentication/Controller"
	"foodr/configuration"
	misc_models "foodr/misc/models"
	user_models "foodr/user/models"

	"github.com/gin-gonic/gin"

	"golang.org/x/crypto/bcrypt"
	"foodr/user"
)

func GetAllUsers(c *gin.Context) {
	var users []user_models.UserDTO
	configuration.DB.Find(&users)

	c.JSON(http.StatusOK, gin.H{"items": users})
}

func CreateUser(c *gin.Context) {
	// Validate input
	var input user_models.UserDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error ": err.Error()})
		return
	}

	// Input validation, fullname cannot be empty
	if len(input.Fullname) <= 0 {
		c.JSON(
			http.StatusBadRequest,
			user.FullnameCannotBeEmpty())
		return
	}

	// Input validation, username cannot be empty
	if len(input.Username) <= 0 {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": misc_models.NewError(
				"Username cannot be null or empty",
				misc_models.UsernameCannotBeNullOrEmpty,
				true,
			)})
		return
	}

	// Check if email exists
	if err := configuration.DB.Where("email = ?", input.Email).First(&input).Error; err == nil {
		c.JSON(
			http.StatusConflict, 
			user.EmailAlreadyExist())
		return
	}

	// Check if username exists
	if err := configuration.DB.Where("username = ?", input.Username).First(&input).Error; err == nil {
		c.JSON(http.StatusConflict,
			user.UsernameAlreadyExist())
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 8)

	input.Password = string(hashed)

	configuration.DB.Create(&input)
	c.JSON(http.StatusOK, gin.H{"item": input})
}

func Login(c *gin.Context) {
	var u user_models.UserDTO

	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	pwd := u.Password

	hashed, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 8)

	hashErr := bcrypt.CompareHashAndPassword(hashed, []byte(u.Password))
	if hashErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error hashing password"})
		fmt.Println("Hash error")
		fmt.Println(hashErr)
		return
	}

	userNotFoundErr := configuration.DB.Where("username = ?", u.Username).First(&u)
	wrongPwdErr := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pwd))

	if userNotFoundErr != nil && wrongPwdErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or password does not match"})
		return
	}

	configuration.GetRedisConnection()
	ts, err := token_auth.CreateToken(u.ID.String())
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	saveErr := token_auth.CreateAuth(u.ID.String(), ts)
	if saveErr != nil {
		c.JSON(http.StatusUnprocessableEntity, saveErr.Error())
	}
	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}
	c.JSON(http.StatusOK, gin.H{"tokens": tokens, "user": u})
}

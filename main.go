package main

import (
	token_auth "foodr/Authentication/Controller"
	"foodr/configuration"
	"log"
	"os"

	recipe_controller "foodr/Recipe/Controller"
	user_controller "foodr/user/controller"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
)

var router = gin.Default()
var client *redis.Client

func init() {
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	client = redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
}

func main() {
	configuration.ConnectDatabase()

	router.GET("/users", user_controller.GetAllUsers)
	router.POST("/users", user_controller.CreateUser)
	router.POST("/users/login", user_controller.Login)

	router.POST("/recipe", recipe_controller.CreateRecipe)
	router.GET("/recipe/:id", recipe_controller.GetRecipe)

	router.POST("/token/refresh", token_auth.Refresh)

	log.Fatal(router.Run(":8080"))
}

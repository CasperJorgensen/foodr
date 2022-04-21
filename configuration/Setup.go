package configuration

import (
	"os"

	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	user_models "foodr/user/models"
	recipe_models "foodr/Recipe/Models"
)

var DB *gorm.DB

func ConnectDatabase() {
	database, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=postgres sslmode=disable password=postgres")

	if err != nil {
		panic("Failed to connect")
	}

	database.AutoMigrate(&user_models.UserDTO{})
	database.AutoMigrate(&recipe_models.RecipeDTO{})

	DB = database
}

func GetRedisConnection() *redis.Client {
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	client := redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
	})
	return client
}
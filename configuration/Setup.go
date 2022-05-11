package configuration

import (
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/go-redis/redis/v7"
	// "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	recipe_models "foodr/Recipe/Models"
	user_models "foodr/user/models"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// database, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=postgres sslmode=disable password=postgres")

	if err != nil {
		panic("Failed to connect")
	}

	db.AutoMigrate(&user_models.UserDTO{})
	db.AutoMigrate(&recipe_models.RecipeDTO{})

	DB = db
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

func Paginate1(value interface{}, pagination *Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(value).Count(&totalRows)

	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}

func Paginate2(r *http.Request) func(db *gorm.DB) *gorm.DB {
	var pg Pagination
	return func(db *gorm.DB) *gorm.DB {
		q := r.URL.Query()
		page, _ := strconv.Atoi(q.Get("page"))
		if page == 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(q.Get("page_size"))
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		return db.Offset(pg.GetLimit()).Limit(pageSize)
	}
}

func Paginate3(r *http.Request) {
	// q := r.URL.Query()
	// page, _ := strconv.Atoi(q.Get("page"))
	// if page <= 0 {
	// 	page = 1
	// }

	// pageSize, _ := strconv.Atoi(q.Get("per_page"))
	// sort := q.Get("sort")
	// var sortWithDirection string
	// if sort != "" {
	// 	direction := q.Get("sortDesc")
	// 	if direction != "" {
	// 		if direction == "true" {
	// 			sortWithDirection = sort + " desc"
	// 		} else if direction == "false" {
	// 			sortWithDirection = sort + " asc"
	// 		}
	// 	}
	// }

	// switch {
	// case pageSize <= 0:
	// 	pageSize = 10
	// }
}

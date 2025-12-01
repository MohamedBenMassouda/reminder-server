package initializers

import (
	"log"
	"os"
	"reminder-server/internal/handlers"
	"reminder-server/internal/models"
	"reminder-server/internal/services"

	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Config struct {
	Port string
}

type Initializers struct {
	CategoryHandler *handlers.CategoryHandler
	ReminderHandler *handlers.ReminderHandler
	UserHandler     *handlers.UserHandler
	Config          *Config
}

var (
	DB *gorm.DB
)

func LoadEnv(files ...string) {
	err := godotenv.Load(files...)

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GetDBString() string {
	url := os.Getenv("DATABASE_URL")

	if url == "" {
		log.Fatal("DATABASE_URL not found")
	}

	authToken := os.Getenv("TURSO_AUTH_TOKEN")

	if authToken == "" {
		log.Fatal("TURSO_AUTH_TOKEN not found")
	}

	return url + "?authToken=" + authToken
}

func ConnectDB() {
	db, err := gorm.Open(sqlite.New(sqlite.Config{
		DriverName: "libsql",
		DSN:        GetDBString(),
	}), &gorm.Config{})

	if err != nil {
		log.Fatal("Error connecting to database")
	}

	DB = db
}

func GetDB() *gorm.DB {
	return DB
}

func NewInitializers() *Initializers {
	LoadEnv()
	ConnectDB()

	categoryService := services.NewCategoryService(DB)
	reminderService := services.NewReminderService(DB)
	userService := services.NewUserService(DB)

	seedCategories(categoryService)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	config := &Config{
		Port: port,
	}

	return &Initializers{
		CategoryHandler: handlers.NewCategoryHandler(categoryService),
		ReminderHandler: handlers.NewReminderHandler(reminderService),
		UserHandler:     handlers.NewUserHandler(userService),
		Config:          config,
	}
}

func seedCategories(categoryService *services.CategoryService) {
	// Check if the categories are already seeded
	dbCategories, err := categoryService.List(1)

	if err == nil && len(dbCategories) > 0 {
		return
	}

	var categories = []models.Category{
		{
			Name:   "Work",
			Icon:   "💼",
			Color:  "#0077B6",
			UserID: 1,
		},
		{
			Name:   "Personal",
			Icon:   "👤",
			Color:  "#FFA500",
			UserID: 1,
		},
		{
			Name:   "Shopping",
			Icon:   "🛒",
			Color:  "#C0392B",
			UserID: 1,
		},
		{
			Name:   "Education",
			Icon:   "📚",
			Color:  "#8E44AD",
			UserID: 1,
		},
	}

	err = categoryService.CreateBulk(categories)

	if err != nil {
		log.Fatal("Error seeding categories")
	}
}

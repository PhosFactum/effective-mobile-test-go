package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/PhosFactum/effective-mobile-test-go/internal/config"
	"github.com/PhosFactum/effective-mobile-test-go/internal/database"
	"github.com/PhosFactum/effective-mobile-test-go/internal/handlers"
	"github.com/PhosFactum/effective-mobile-test-go/internal/repository"

	_ "github.com/PhosFactum/effective-mobile-test-go/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Загрузка .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := config.Load()

	// Подключение к БД
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Применение миграций
	if err := database.RunMigrations(db, cfg); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Инициализация репозитория и хендлеров
	repo := repository.NewSubscriptionRepository(db)
	h := handlers.NewSubscriptionHandler(repo)

	// Настройка роутера
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// CRUDL
	api := r.Group("/api/v1")
	{
		api.POST("/subscriptions", h.Create)
		api.GET("/subscriptions", h.List)
		api.GET("/subscriptions/:id", h.GetByID)
		api.PUT("/subscriptions/:id", h.Update)
		api.DELETE("/subscriptions/:id", h.Delete)
		api.GET("/subscriptions/total-cost", h.TotalCost)
	}

	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

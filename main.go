package main

import (
	"fmt"
	"log"

	"github.com/ricirt/webhook-automation/internal/config"
	"github.com/ricirt/webhook-automation/internal/handler"
	"github.com/ricirt/webhook-automation/internal/repository"
	"github.com/ricirt/webhook-automation/internal/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Insider Message Sending API
// @version 1.0
// @description API for automatic message sending system
// @host localhost:8080
// @BasePath /api/v1

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize repositories
	messageRepo := repository.NewMessageRepository(db)

	// Initialize services
	messageService := service.NewMessageService(messageRepo)

	// Initialize handlers
	messageHandler := handler.NewMessageHandler(messageService)

	// Initialize router
	r := gin.Default()

	// API routes
	api := r.Group("/api/v1")
	{
		messages := api.Group("/messages")
		{
			messages.GET("/sent", messageHandler.GetSentMessages)
			messages.GET("", messageHandler.GetAllMessages)
		}
	}

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	_ "github.com/ricirt/webhook-automation/docs"
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

// @title           Insider Message Sending API
// @version         1.0
// @description     API for automatic message sending system
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @schemes   http https

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

	// Initialize Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
	})

	// Test Redis connection
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Initialize repositories
	messageRepo := repository.NewMessageRepository(db)

	// Initialize services
	messageService := service.NewMessageService(messageRepo, cfg, rdb)

	// Initialize handlers
	messageHandler := handler.NewMessageHandler(messageService)

	// Initialize router
	r := gin.Default()

	// API routes
	api := r.Group("/api/v1")
	{
		messages := api.Group("/messages")
		{
			messages.POST("/start", messageHandler.StartSending)
			messages.POST("/stop", messageHandler.StopSending)
			messages.GET("/sent", messageHandler.GetSentMessages)
		}
	}

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

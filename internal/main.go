package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/warleon/ms4-compliance-service/internal/config"
	"github.com/warleon/ms4-compliance-service/internal/handlers"
	"github.com/warleon/ms4-compliance-service/internal/middleware"
	"github.com/warleon/ms4-compliance-service/internal/repository"
	"github.com/warleon/ms4-compliance-service/internal/repository/rules"
	"github.com/warleon/ms4-compliance-service/internal/service"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/warleon/ms4-compliance-service/docs"
)

// @title Compliance Rules API
// @version 1.0
// @description This is the API documentation for the Compliance Rules service.
// @host localhost:8080
// @BasePath /api/v1
// @schemes http
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found, relying on environment variables")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	db, err := config.NewGormDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	// run auto-migrations
	for _, r := range rules.RuleTables {
		db.AutoMigrate(r)
	}
	for _, r := range repository.RepositoryTables {
		db.AutoMigrate(r)
	}

	repo := repository.NewMySQLRepository(db)
	compService := service.NewComplianceService(repo)
	handler := handlers.NewComplianceHandler(compService)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogger())

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"content": "Hola mundo"})
	})

	api := r.Group("/api/v1")
	{
		api.POST("/validateTransaction", handler.ValidateTransaction)
		api.POST("/rules", handler.CreateRule)
		api.GET("/rules/:id", handler.GetRule)
		api.GET("/rules", handler.ListRules)
		api.PUT("/rules/:id", handler.UpdateRule)
		api.DELETE("/rules/:id", handler.DeleteRule)
	}

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	addr := fmt.Sprintf(":%s", cfg.Port)
	logrus.Infof("starting server on %s", addr)
	r.Run(addr)
}

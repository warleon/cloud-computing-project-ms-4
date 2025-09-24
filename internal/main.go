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
	"github.com/warleon/ms4-compliance-service/internal/service"
)

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
	db.AutoMigrate(&repository.Rule{}, &repository.AuditLog{})

	repo := repository.NewMySQLRepository(db)
	fraudClient := service.NewHTTPFraudClient(cfg.FraudAPIURL)
	compService := service.NewComplianceService(repo, fraudClient)
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
		api.GET("/riskScore/:customerId", handler.GetRiskScore)
		api.POST("/rules", handler.CreateRule)
		api.GET("/rules", handler.ListRules)
	}

	addr := fmt.Sprintf(":%s", cfg.Port)
	logrus.Infof("starting server on %s", addr)
	r.Run(addr)
}

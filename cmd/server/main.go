package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/silvioubaldino/ilia-wallet/internal/adapter/http/middleware"
	"github.com/silvioubaldino/ilia-wallet/internal/infrastructure/bootstrap"
	"github.com/silvioubaldino/ilia-wallet/internal/infrastructure/config"
	"github.com/silvioubaldino/ilia-wallet/internal/infrastructure/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.NewPostgres(cfg.DSN(), cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	auth := router.Group("/")
	auth.Use(middleware.Auth(cfg.JWTSecret))

	bootstrap.SetupTransaction(db, auth)

	log.Printf("server starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/supel2nova/welfare-registration/backend/internal/config"
	"github.com/supel2nova/welfare-registration/backend/internal/repository"
)

func main() {
	_ = godotenv.Load("../.env")
	cfg := config.Load()

	ctx := context.Background()
	pool, err := repository.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer pool.Close()

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		if err := pool.Ping(c); err != nil {
			c.JSON(503, gin.H{"status": "db down"})
			return
		}
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Printf("listening on :%s", cfg.HTTPPort)
	_ = r.Run(":" + cfg.HTTPPort)
}

package main

import (
	"languages-api/internal/config"
	"languages-api/internal/models"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()

	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	db, err := config.NewDatabase(cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database connected successfully")

	db.AutoMigrate(&models.Language{})

	var router *gin.Engine = gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Programming Languages API is running",
			"status":  "success",
		})
	})

	router.Run(":" + cfg.Port)
}

package main

import (
	"languages-api/internal/config"
	"languages-api/internal/handlers"
	"languages-api/internal/models"
	"languages-api/internal/repository"
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

	repo := repository.NewLanguageRepository(db)
	handler := handlers.NewLanguageHandler(repo)

	var router *gin.Engine = gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Programming Languages API is running",
			"status":  "success",
		})
	})

	api := router.Group("/api/v1")
	{
		// READ
		api.GET("/languages/:id", handler.GetLanguageByID)
		api.GET("/languages", handler.GetLanguages)
		// CREATE
		api.POST("/languages", handler.CreateLanguage)
		api.POST("/languages/batch", handler.CreateLanguages)
		// UPDATE
		api.PATCH("/languages/:id", handler.UpdateLanguage)
		// DELETE
		api.DELETE("/languages/:id", handler.DeleteLanguage)
		api.POST("/languages/:id/restore", handler.RestoreLanguage)
	}

	router.Run(":" + cfg.Port)
}

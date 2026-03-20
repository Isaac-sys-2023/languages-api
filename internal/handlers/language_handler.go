package handlers

import (
	"errors"
	"languages-api/internal/models"
	"languages-api/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LanguageHandler struct {
	repo repository.LanguageRepository
}

func NewLanguageHandler(r repository.LanguageRepository) *LanguageHandler {
	return &LanguageHandler{
		repo: r,
	}
}

func (h *LanguageHandler) GetLanguages(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	year, _ := strconv.Atoi(c.Query("release_year"))

	filters := repository.LanguageFilters{
		Name:        c.Query("name"),
		Creator:     c.Query("creator"),
		ReleaseYear: year,
	}

	languages, total, err := h.repo.FindWithFilters(c.Request.Context(), filters, page, pageSize)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching languages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": languages,
		"meta": gin.H{
			"total_records": total,
			"current_page":  page,
			"page_size":     pageSize,
		},
	})
}

func (h *LanguageHandler) GetLanguageByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	language, err := h.repo.FindByID(c.Request.Context(), uint(id))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if language == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Language not found"})
		return
	}

	c.JSON(http.StatusOK, language)
}

func (h *LanguageHandler) CreateLanguage(c *gin.Context) {
	var lang models.Language

	// "Bind" del JSON al struct (Validación de tipos automática)
	if err := c.ShouldBindJSON(&lang); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format", "details": err.Error()})
		return
	}

	// Llamar al repositorio pasándole el puntero del lenguaje
	if err := h.repo.Create(c.Request.Context(), &lang); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create language"})
		return
	}

	// Respuesta 201 Created con el objeto ya guardado (incluyendo su ID de DB)
	c.JSON(http.StatusCreated, lang)
}

func (h *LanguageHandler) CreateLanguages(c *gin.Context) {
	var languages []models.Language

	// Intentar bindear a un slice/array de JSON
	if err := c.ShouldBindJSON(&languages); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Expected a list of languages"})
		return
	}

	// Llamar al método CreateBatch que ya definimos en el repo
	if err := h.repo.CreateBatch(c.Request.Context(), languages); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store batch"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Batch processed successfully",
		"count":   len(languages),
	})
}

func (h *LanguageHandler) UpdateLanguage(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	err := h.repo.Update(c.Request.Context(), uint(id), updates)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Language not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Language updated successfully"})
}

func (h *LanguageHandler) DeleteLanguage(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.repo.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent) // 204: Éxito pero sin cuerpo de respuesta
}

func (h *LanguageHandler) RestoreLanguage(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.repo.Restore(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Restore failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Language restored"})
}

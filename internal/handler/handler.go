package handler

import (
	"github.com/DebroyeAntoine/go_link_vault/internal/db"
	"github.com/DebroyeAntoine/go_link_vault/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateLinkHandler(c *gin.Context) {
	var link models.Link
	if err := c.ShouldBindJSON(&link); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Create(&link).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save the link"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": link.ID, "url": link.URL, "title": link.Title, "tags": link.Tags})
}

package handler

import (
	"encoding/json"
	"net/http"

	"github.com/DebroyeAntoine/go_link_vault/internal/auth"
	"github.com/DebroyeAntoine/go_link_vault/internal/db"
	"github.com/DebroyeAntoine/go_link_vault/internal/dto"
	"github.com/DebroyeAntoine/go_link_vault/internal/logger"
	"github.com/DebroyeAntoine/go_link_vault/internal/models"
	"github.com/DebroyeAntoine/go_link_vault/internal/scraper"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

func CreateLinkHandler(c *gin.Context) {
	// Extraire l'email du token
	userEmail, exists := c.Get("userEmail")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Trouver l'utilisateur dans la DB
	var user models.User
	if err := db.DB.Where("email = ?", userEmail).First(&user).Error; err != nil {
		ErrorResponse(c, http.StatusUnauthorized, "User not found")
		return
	}

	var input dto.CreateLinkDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	tagsJSON, _ := json.Marshal(input.Tags)
	link := models.Link{
		URL:    input.URL,
		Title:  input.Title,
		Tags:   datatypes.JSON(tagsJSON),
		UserID: user.ID,
	}
	if err := db.DB.Create(&link).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "could not save the link")
		return
	}

	go func(linkID uint, url string) {
		metadata, err := scraper.FetchMetadata(url)
		if err != nil {
			logger.ErrorLogger.Println("Failed to fetch metadata:", err)
			return
		}

		db.DB.Model(&models.Link{}).Where("id = ?", linkID).Updates(models.Link{
			Description: metadata.Description,
			Image:       metadata.Image,
		})
	}(link.ID, link.URL)

	SuccessResponse(c, http.StatusCreated, gin.H{
		"url":   link.URL,
		"title": link.Title,
		"tags":  link.Tags,
	})
}

// Register handler
func RegisterUserHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}
	user.Password = hashedPassword

	// Save user to DB
	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	// Generate JWT token
	token, err := auth.CreateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating JWT"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token})
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginUserHandler(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dbUser models.User
	if err := db.DB.Where("email = ?", input.Email).First(&dbUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !auth.CheckPasswordHash(input.Password, dbUser.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := auth.CreateToken(dbUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating JWT"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func GetLinksHandler(c *gin.Context) {
	// Vérification du token JWT et récupération de l'ID de l'utilisateur
	claims, err := auth.ValidateToken(c)
	if err != nil {
		ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Récupération de l'email directement depuis les claims
	email := claims.Issuer

	var user models.User
	if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		ErrorResponse(c, http.StatusUnauthorized, "User not found")
		return
	}

	// Récupération de tous les liens de l'utilisateur
	var links []models.Link /* No preload because useless and risky to return User */
	if err := db.DB.Where("user_id = ?", user.ID).Find(&links).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Could not fetch links")
		return
	}

	SuccessResponse(c, http.StatusOK, links)
}

func UpdateLinkHandler(c *gin.Context) {
	claims, err := auth.ValidateToken(c)
	if err != nil {
		ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// On récupère l'utilisateur
	var user models.User
	if err := db.DB.Where("email = ?", claims.Issuer).First(&user).Error; err != nil {
		ErrorResponse(c, http.StatusUnauthorized, "User not found")
		return
	}

	// On récupère l'ID du lien depuis l'URL
	linkID := c.Param("id")
	var link models.Link
	if err := db.DB.Where("id = ? AND user_id = ?", linkID, user.ID).First(&link).Error; err != nil {
		ErrorResponse(c, http.StatusNotFound, "Link not found")
		return
	}

	var input dto.UpdateLinkDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid input")
		return
	}

	// Mise à jour des champs modifiables
	if input.URL != nil {
		link.URL = *input.URL
	}
	if input.Title != nil {
		link.Title = *input.Title
	}
	if input.Tags != nil {
		tagsJSON, _ := json.Marshal(*input.Tags)
		link.Tags = tagsJSON
	}

	if err := db.DB.Save(&link).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Could not update the link")
		return
	}

	SuccessResponse(c, http.StatusOK, gin.H{
		"id":    link.ID,
		"url":   link.URL,
		"title": link.Title,
		"tags":  link.Tags,
	})
}

func DeleteLinkHandler(c *gin.Context) {
	claims, err := auth.ValidateToken(c)
	if err != nil {
		ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	email := claims.Issuer
	var user models.User
	if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		ErrorResponse(c, http.StatusUnauthorized, "User not found")
		return
	}

	linkID := c.Param("id")
	var link models.Link
	if err := db.DB.Where("id = ? AND user_id = ?", linkID, user.ID).First(&link).Error; err != nil {
		ErrorResponse(c, http.StatusNotFound, "Link not found")
		return
	}

	if err := db.DB.Delete(&link).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Could not delete link")
		return
	}

	SuccessResponse(c, http.StatusOK, gin.H{"message": "Link deleted successfully"})
}

func GetLinkHandler(c *gin.Context) {
	claims, err := auth.ValidateToken(c)
	if err != nil {
		ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// On récupère l'utilisateur
	var user models.User
	if err := db.DB.Where("email = ?", claims.Issuer).First(&user).Error; err != nil {
		ErrorResponse(c, http.StatusUnauthorized, "User not found")
		return
	}

	// On récupère l'ID du lien depuis l'URL
	linkID := c.Param("id")
	var link models.Link
	if err := db.DB.Where("id = ? AND user_id = ?", linkID, user.ID).First(&link).Error; err != nil {
		ErrorResponse(c, http.StatusNotFound, "Link not found")
		return
	}
	SuccessResponse(c, http.StatusOK, link)
}

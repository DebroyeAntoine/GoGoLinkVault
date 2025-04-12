package handler

import (
	"net/http"

	"github.com/DebroyeAntoine/go_link_vault/internal/auth"
	"github.com/DebroyeAntoine/go_link_vault/internal/db"
	"github.com/DebroyeAntoine/go_link_vault/internal/models"
	"github.com/gin-gonic/gin"
)

// Link Handler
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

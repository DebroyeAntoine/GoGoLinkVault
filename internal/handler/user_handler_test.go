package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DebroyeAntoine/go_link_vault/internal/auth"
	"github.com/DebroyeAntoine/go_link_vault/internal/db"
	"github.com/DebroyeAntoine/go_link_vault/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	db.SetupTestDB()
	router := gin.Default()
	router.POST("/register", RegisterUserHandler)

	payload := map[string]string{
		"email":    "newuser@example.com",
		"password": "securepass123",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, 201, resp.Code)

	var response map[string]string
	json.Unmarshal(resp.Body.Bytes(), &response)

	assert.NotEmpty(t, response["token"])
}

func TestLoginUser(t *testing.T) {
	db.SetupTestDB()

	hashedpwd, _ := auth.HashPassword("hashedpassword123")

	user := models.User{
		Email:    "test@example.com",
		Password: hashedpwd,
	}

	err := db.DB.Create(&user).Error
	assert.NoError(t, err)

	router := gin.Default()
	router.POST("/login", LoginUserHandler)

	payload := map[string]string{
		"email":    "test@example.com",
		"password": "hashedpassword123",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)

	var response map[string]string
	json.Unmarshal(resp.Body.Bytes(), &response)

	assert.NotEmpty(t, response["token"])
}

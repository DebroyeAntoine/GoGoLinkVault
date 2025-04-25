package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DebroyeAntoine/go_link_vault/internal/auth"
	"github.com/DebroyeAntoine/go_link_vault/internal/db"
	"github.com/DebroyeAntoine/go_link_vault/internal/logger"
	"github.com/DebroyeAntoine/go_link_vault/internal/middleware"
	"github.com/DebroyeAntoine/go_link_vault/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

func toJSON(tags []string) datatypes.JSON {
	jsonBytes, _ := json.Marshal(tags)
	return datatypes.JSON(jsonBytes)
}

type ResponseData[T any] struct {
	Success bool `json:"success"`
	Data    T    `json:"data"`
}

func TestCreateLink(t *testing.T) {
	logger.InitLogger()
	// Configurez la base de données de test
	db.SetupTestDB()

	// Créez un utilisateur et générez un token JWT
	hashedpwd, _ := auth.HashPassword("hashedpassword123")
	user := models.User{
		Email:    "test@example.com",
		Password: hashedpwd,
	}

	err := db.DB.Create(&user).Error
	if err != nil {
		t.Fatalf("Error creating user: %v", err)
	}

	token, err := auth.CreateToken(user)
	if err != nil {
		t.Fatalf("Error creating JWT token: %v", err)
	}

	// Configurez le routeur Gin
	r := gin.Default()

	// Protéger la route avec le middleware AuthRequired
	r.POST("/links", middleware.AuthRequired(), CreateLinkHandler)

	// Créez la charge utile pour la requête
	payload := map[string]interface{}{
		"url":   "https://go.dev",
		"title": "The Go Programming Language",
		"tags":  []string{"go", "programming"},
	}
	jsonBody, _ := json.Marshal(payload)

	// Créez une requête HTTP de test avec un en-tête Authorization
	req, _ := http.NewRequest("POST", "/links", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()

	// Exécutez la requête
	r.ServeHTTP(resp, req)

	// FIXME remove sleep with mock test
	time.Sleep(500 * time.Millisecond)

	// Vérifiez les résultats
	assert.Equal(t, http.StatusCreated, resp.Code)

	var link models.Link
	db.DB.Last(&link)

	var tags []string
	err = json.Unmarshal(link.Tags, &tags)
	if err != nil {
		t.Fatalf("Failed to unmarshal tags: %v", err)
	}

	assert.Equal(t, "https://go.dev", link.URL)
	assert.Equal(t, "The Go Programming Language", link.Title)
	assert.ElementsMatch(t, []string{"go", "programming"}, tags)
}

func TestGetLinks(t *testing.T) {
	db.SetupTestDB()

	// Création d’un utilisateur et hash du mot de passe
	hashedPwd, _ := auth.HashPassword("testpassword")
	user := models.User{
		Email:    "getlinks@example.com",
		Password: hashedPwd,
	}
	err := db.DB.Create(&user).Error
	assert.NoError(t, err)

	// Création de quelques liens pour cet utilisateur
	links := []models.Link{
		{
			URL:    "https://golang.org",
			Title:  "Golang",
			Tags:   toJSON([]string{"go", "lang"}),
			UserID: user.ID,
		},
		{
			URL:    "https://gin-gonic.com",
			Title:  "Gin Web Framework",
			Tags:   toJSON([]string{"go", "web", "framework"}),
			UserID: user.ID,
		},
	}
	for _, link := range links {
		assert.NoError(t, db.DB.Create(&link).Error)
	}

	// Génère un token JWT
	token, err := auth.CreateToken(user)
	assert.NoError(t, err)

	// Setup du routeur avec middleware
	r := gin.Default()
	r.GET("/links", middleware.AuthRequired(), GetLinksHandler)

	// Requête GET avec header Authorization
	req, _ := http.NewRequest("GET", "/links", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	// Désérialisation de la réponse JSON dans un objet avec une clé `data`
	var jsonResponse ResponseData[[]models.Link]
	err = json.Unmarshal(resp.Body.Bytes(), &jsonResponse)
	assert.NoError(t, err)

	// Vérifie que nous avons deux liens dans la réponse
	assert.Len(t, jsonResponse.Data, 2)

	// Vérifie les détails des liens
	assert.Equal(t, "https://golang.org", jsonResponse.Data[0].URL)
	assert.Equal(t, "Golang", jsonResponse.Data[0].Title)

	assert.Equal(t, "https://gin-gonic.com", jsonResponse.Data[1].URL)
	assert.Equal(t, "Gin Web Framework", jsonResponse.Data[1].Title)
}

func TestCreateLinkValidation(t *testing.T) {
	// setup user, token, route...
	db.SetupTestDB()

	// Création d’un utilisateur et hash du mot de passe
	hashedPwd, _ := auth.HashPassword("testpassword")
	user := models.User{
		Email:    "getlinks@example.com",
		Password: hashedPwd,
	}
	err := db.DB.Create(&user).Error
	assert.NoError(t, err)

	token, err := auth.CreateToken(user)
	assert.NoError(t, err)

	// Setup du routeur avec middleware
	r := gin.Default()
	r.POST("/links", middleware.AuthRequired(), CreateLinkHandler)

	payload := `{"title": "Missing URL"}`
	req, _ := http.NewRequest("POST", "/links", strings.NewReader(payload))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "URL")
}

func TestUpdateLink(t *testing.T) {
	db.SetupTestDB()

	hashedPwd, _ := auth.HashPassword("password")
	user := models.User{
		Email:    "update@example.com",
		Password: hashedPwd,
	}
	assert.NoError(t, db.DB.Create(&user).Error)

	link := models.Link{
		URL:    "https://old-url.com",
		Title:  "Old Title",
		Tags:   toJSON([]string{"old", "tag"}),
		UserID: user.ID,
	}
	assert.NoError(t, db.DB.Create(&link).Error)

	token, err := auth.CreateToken(user)
	assert.NoError(t, err)

	r := gin.Default()
	r.PUT("/links/:id", middleware.AuthRequired(), UpdateLinkHandler)

	newTitle := "New Title"
	newTags := []string{"new", "cool"}
	updateData := map[string]interface{}{
		"title": newTitle,
		"tags":  newTags,
	}
	body, _ := json.Marshal(updateData)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/links/%d", link.ID), bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, newTitle, data["title"])
}

func TestGetLinkByID(t *testing.T) {
	db.SetupTestDB()

	// Création de l'utilisateur
	hashedPwd, _ := auth.HashPassword("testpassword")
	user := models.User{
		Email:    "getone@example.com",
		Password: hashedPwd,
	}
	assert.NoError(t, db.DB.Create(&user).Error)

	// Création d’un lien
	link := models.Link{
		URL:    "https://example.com",
		Title:  "Example",
		Tags:   toJSON([]string{"test"}),
		UserID: user.ID,
	}
	assert.NoError(t, db.DB.Create(&link).Error)

	// Token JWT
	token, _ := auth.CreateToken(user)

	// Setup route
	r := gin.Default()
	r.GET("/links/:id", middleware.AuthRequired(), GetLinkHandler)

	// Requête
	req, _ := http.NewRequest("GET", fmt.Sprintf("/links/%d", link.ID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response ResponseData[models.Link]
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "https://example.com", response.Data.URL)
}

func TestDeleteLink(t *testing.T) {
	db.SetupTestDB()

	// Création de l'utilisateur
	hashedPwd, _ := auth.HashPassword("deletepass")
	user := models.User{
		Email:    "delete@example.com",
		Password: hashedPwd,
	}
	assert.NoError(t, db.DB.Create(&user).Error)

	// Création d’un lien
	link := models.Link{
		URL:    "https://delete.me",
		Title:  "To delete",
		Tags:   toJSON([]string{"remove"}),
		UserID: user.ID,
	}
	assert.NoError(t, db.DB.Create(&link).Error)

	// Token JWT
	token, _ := auth.CreateToken(user)

	// Setup route
	r := gin.Default()
	r.DELETE("/links/:id", middleware.AuthRequired(), DeleteLinkHandler)

	// Requête
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/links/%d", link.ID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	// Vérifie que le lien a bien été supprimé
	var deleted models.Link
	err := db.DB.First(&deleted, link.ID).Error
	assert.Error(t, err) // devrait être une "record not found"
}

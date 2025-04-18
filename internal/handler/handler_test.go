package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DebroyeAntoine/go_link_vault/internal/auth"
	"github.com/DebroyeAntoine/go_link_vault/internal/db"
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

func TestCreateLink(t *testing.T) {
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
	var jsonResponse struct {
		Success bool          `json:"success"`
		Data    []models.Link `json:"data"`
	}

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

package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DebroyeAntoine/go_link_vault/internal/db"
	"github.com/DebroyeAntoine/go_link_vault/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateLink(t *testing.T) {
	// Configurez la base de données de test
	db.SetupTestDB()

	// Configurez le routeur Gin
	r := gin.Default()
	r.POST("/links", CreateLinkHandler)

	// Créez la charge utile pour la requête
	payload := map[string]interface{}{
		"url":   "https://go.dev",
		"title": "The Go Programming Language",
		"tags":  []string{"go", "programming"},
	}
	jsonBody, _ := json.Marshal(payload)

	// Créez une requête HTTP de test
	req, _ := http.NewRequest("POST", "/links", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	// Exécutez la requête
	r.ServeHTTP(resp, req)

	// Vérifiez les résultats
	assert.Equal(t, http.StatusCreated, resp.Code)

	var link models.Link
	db.DB.Last(&link)

	var tags []string
	err := json.Unmarshal(link.Tags, &tags)
	if err != nil {
		t.Fatalf("Failed to unmarshal tags: %v", err)
	}

	assert.Equal(t, "https://go.dev", link.URL)
	assert.Equal(t, "The Go Programming Language", link.Title)
	assert.ElementsMatch(t, []string{"go", "programming"}, tags)
}

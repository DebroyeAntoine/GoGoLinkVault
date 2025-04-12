package middleware

import (
	"github.com/DebroyeAntoine/go_link_vault/internal/auth"
	"github.com/DebroyeAntoine/go_link_vault/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthRequired(t *testing.T) {
	// Simuler un router Gin avec le middleware
	router := gin.Default()

	// Route protégée
	router.GET("/protected", AuthRequired(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "You are authorized!"})
	})

	t.Run("Without Token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/protected", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		// Vérifier que la réponse est une erreur 401 Unauthorized
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		assert.JSONEq(t, `{"error": "Missing or invalid Authorization header"}`, resp.Body.String())
	})

	t.Run("With Invalid Token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalidToken")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		// Vérifier que la réponse est une erreur 401 Unauthorized
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		assert.JSONEq(t, `{"error": "Invalid or expired token"}`, resp.Body.String())
	})

	t.Run("With Valid Token", func(t *testing.T) {
		// Créer un token JWT valide
		token, _ := auth.CreateToken(models.User{
			Email: "test@example.com",
		})

		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		// Vérifier que la réponse est OK et le message est correct
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, `{"message": "You are authorized!"}`, resp.Body.String())
	})
}

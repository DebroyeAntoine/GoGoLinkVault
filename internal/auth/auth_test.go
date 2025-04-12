package auth

import (
	"testing"
	"time"

	"github.com/DebroyeAntoine/go_link_vault/internal/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestCreateToken(t *testing.T) {
	// Simuler un utilisateur
	user := models.User{
		Email: "testuser@example.com",
	}

	// Créer un token
	token, err := CreateToken(user)
	if err != nil {
		t.Fatalf("Error creating token: %v", err)
	}

	// Décoder le token pour vérifier les informations
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Vérifie que l'algorithme est HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			t.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil {
		t.Fatalf("Error parsing token: %v", err)
	}

	// Vérifier que les claims sont valides
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		t.Fatalf("Invalid token")
	}

	// Vérifier que l'email est correctement assigné à l'issuer
	assert.Equal(t, user.Email, claims["iss"])

	// Vérifier que le token n'a pas expiré
	expiration := claims["exp"].(float64)
	assert.Greater(t, int64(expiration), time.Now().Unix(), "Token expired")
}

func TestHashPassword(t *testing.T) {
	password := "password123"

	// Créer un hash pour le mot de passe
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}

	// Vérifier que le hash n'est pas égal au mot de passe d'origine
	assert.NotEqual(t, password, hash, "Password hash should not be equal to plain password")

	// Vérifier que le hash est valide avec la fonction CheckPasswordHash
	isValid := CheckPasswordHash(password, hash)
	assert.True(t, isValid, "Password should match hash")
}

func TestCheckPasswordHash(t *testing.T) {
	// Mot de passe de test
	password := "password123"

	// Créer un hash pour le mot de passe
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}

	// Vérifier la correspondance correcte
	isValid := CheckPasswordHash(password, hash)
	assert.True(t, isValid, "Password should match hash")

	// Vérifier une mauvaise correspondance
	incorrectPassword := "wrongpassword"
	isInvalid := CheckPasswordHash(incorrectPassword, hash)
	assert.False(t, isInvalid, "Password should not match hash")
}


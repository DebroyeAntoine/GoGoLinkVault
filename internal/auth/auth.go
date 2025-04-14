package auth

import (
	"time"

	"github.com/DebroyeAntoine/go_link_vault/internal/models"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	"errors"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var jwtKey []byte

func init() {
	key := os.Getenv("JWT_SECRET_KEY")
	if key == "" {
		key = "dev_default_secret" // à utiliser seulement en dev
	}
	jwtKey = []byte(key)
}

func JwtKey() []byte {
	return jwtKey
}

// Create JWT Token
func CreateToken(user models.User) (string, error) {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		Issuer:    user.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Hash password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Check if password is correct
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ValidateToken(c *gin.Context) (*jwt.StandardClaims, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, errors.New("missing authorization header")
	}

	// Extrait le token du header
	tokenString := strings.Split(authHeader, " ")[1]

	// Parse le token
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Récupère les claims (informations utilisateur)
	claims, ok := token.Claims.(*jwt.StandardClaims)
	c.Set("userEmail", claims.Issuer)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

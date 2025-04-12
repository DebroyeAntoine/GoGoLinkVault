package db

import (
	"testing"

	"github.com/DebroyeAntoine/go_link_vault/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestUserModel(t *testing.T) {
	SetupTestDB()

	user := models.User{
		Email:    "test@example.com",
		Password: "hashedpassword123",
		JWTToken: "token123",
	}

	err := DB.Create(&user).Error
	assert.NoError(t, err)

	var found models.User
	err = DB.Where("email = ?", "test@example.com").First(&found).Error
	assert.NoError(t, err)
	assert.Equal(t, user.Email, found.Email)
	assert.Equal(t, user.Password, found.Password)
	assert.Equal(t, user.JWTToken, found.JWTToken)
}


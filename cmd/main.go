package main

import (
	"github.com/DebroyeAntoine/go_link_vault/internal/db"
	"github.com/DebroyeAntoine/go_link_vault/internal/link"
	"github.com/gin-gonic/gin"
)

func main() {
	// Connexion DB
	db.Connect()

	r := gin.Default()

	r.POST("/links", link.CreateLinkHandler)

	r.Run(":8080")
}

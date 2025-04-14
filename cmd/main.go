package main

import (
	"github.com/DebroyeAntoine/go_link_vault/internal/db"
	"github.com/DebroyeAntoine/go_link_vault/internal/handler"
	"github.com/DebroyeAntoine/go_link_vault/internal/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	// Connexion DB
	db.Connect()

	r := gin.Default()

	r.POST("/register", handler.RegisterUserHandler)
	r.POST("/links", handler.LoginUserHandler)
	r.POST("/links", middleware.AuthRequired(), handler.CreateLinkHandler)
	r.GET("/links", middleware.AuthRequired(), handler.GetLinksHandler)

	r.Run(":8080")
}

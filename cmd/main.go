package main

import (
	"time"

	"github.com/DebroyeAntoine/go_link_vault/internal/db"
	"github.com/DebroyeAntoine/go_link_vault/internal/handler"
	"github.com/DebroyeAntoine/go_link_vault/internal/logger"
	"github.com/DebroyeAntoine/go_link_vault/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	logger.InitLogger()
	// Connexion DB
	db.Connect()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/register", handler.RegisterUserHandler)
	r.POST("/login", handler.LoginUserHandler)
	r.POST("/links", middleware.AuthRequired(), handler.CreateLinkHandler)
	r.GET("/links", middleware.AuthRequired(), handler.GetLinksHandler)
	r.PUT("/link/:id", middleware.AuthRequired(), handler.UpdateLinkHandler)
	r.DELETE("link/:id", middleware.AuthRequired(), handler.DeleteLinkHandler)
	r.GET("link/:id", middleware.AuthRequired(), handler.GetLinkHandler)

	r.Run(":8080")
}

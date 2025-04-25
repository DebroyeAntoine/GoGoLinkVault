package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Tags est un type personnalisé qui implémente l'interface gorm.DB pour être stocké comme JSON
type Tags []string

// Link représente un lien avec son URL, son titre et ses tags
type Link struct {
	gorm.Model
	URL         string         `json:"url" binding:"required,url"`
	Title       string         `json:"title" binding:"required"`
	Tags        datatypes.JSON `json:"tags"`
	UserID      uint           `json:"-"` // Clé étrangère
	User        User           `gorm:"foreignKey:UserID" json:"-"`
	Description string         `json:"description,omitempty"`
	Image       string         `json:"image,omitempty"`
}

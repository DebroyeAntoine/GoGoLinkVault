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
	URL   string         `json:"url"`
	Title string         `json:"title"`
	Tags  datatypes.JSON `json:"tags"`
}

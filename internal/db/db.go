package db

import (
	"fmt"
	"log"
	"os"

	"github.com/DebroyeAntoine/go_link_vault/internal/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Fonction pour connecter à la base de données
func Connect() {
	// Charger le fichier .env
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	// Récupérer les variables d'environnement
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	// Construire la chaîne de connexion
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode)

	// Connexion à la base de données
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	// Auto-migrate pour créer les tables si elles n'existent pas
	err = DB.AutoMigrate(&models.User{}, &models.Link{})
	if err != nil {
		log.Fatal("Error migrating database: ", err)
	}
}

// Fonction utilitaire pour configurer une base de données de test
func SetupTestDB() {
	// Charger le fichier .env
	env_err := godotenv.Load("../../.env")
	if env_err != nil {
		log.Fatal("Error loading .env file", env_err)
	}

	// Récupérer les variables d'environnement pour la base de données de test
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	testDBName := os.Getenv("DB_TEST_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	// Construire la chaîne de connexion pour la base de données de test
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, testDBName, port, sslmode)

	// Connexion à la base de données de test
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to test database")
	}

	// Supprimer la table existante si elle existe
	DB.Exec("TRUNCATE TABLE links, users RESTART IDENTITY CASCADE")

	err = DB.AutoMigrate(&models.User{}, &models.Link{})
	if err != nil {
		log.Fatal("Error migrating database: ", err)
	}
}

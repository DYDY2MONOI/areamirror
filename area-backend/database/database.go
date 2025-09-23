package database

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Fichier .env requis mais introuvable:", err)
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")

	if port == "" {
		port = "5432"
	}
	if sslmode == "" {
		sslmode = "disable"
	}


	if host == "" || user == "" || dbname == "" {
		log.Fatal("Variables d'environnement manquantes dans .env: DB_HOST, DB_USER, DB_NAME sont obligatoires")
	}

	dsn := "host=" + host + " user=" + user
	if password != "" {
		dsn += " password=" + password
	}
	dsn += " dbname=" + dbname + " port=" + port + " sslmode=" + sslmode


	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Échec de connexion à la base de données:", err)
	}

	log.Println("Connexion à PostgreSQL réussie!")
}

func GetDB() *gorm.DB {
	return DB
}

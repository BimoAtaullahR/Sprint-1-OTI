package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // Driver PostgreSQL
)

var DB *sql.DB

func ConnectDB() {
	// Membaca dari environment variable Docker. Jika kosong, gunakan default (untuk run lokal).
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "password123"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "mbg_inventory"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Gagal membuka koneksi: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Database tidak merespon: %v", err)
	}

	DB = db
	log.Println("Berhasil terhubung ke database!")
}

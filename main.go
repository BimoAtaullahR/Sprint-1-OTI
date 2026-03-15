package main

import (
	"log"

	"school-vault/config"
	"school-vault/routes"
)

func main() {
	// 1. Konek ke Database
	config.ConnectDB()

	// 2. Setup Routes
	r := routes.SetupRouter()

	// 3. Jalankan Server
	log.Println("Server berjalan di port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}

package main

import (
	"log"

	"github.com/hafizhproject45/Golang-Boilerplate.git/internal/config"
	"github.com/hafizhproject45/Golang-Boilerplate.git/internal/database"
	"github.com/hafizhproject45/Golang-Boilerplate.git/internal/database/seed"
)

func main() {
	db := database.Connect(config.DBHost, config.DBName)

	if err := seed.Run(db); err != nil {
		log.Fatalf("❌ Failed run seeder: %v", err)
	}

	log.Println("✅ Seed Successfully")
}

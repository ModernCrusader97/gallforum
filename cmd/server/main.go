package main

import (
	"log"
	"os"

	"arcalive/internal/repository"
	"arcalive/internal/router"
)

func main() {
	db, err := repository.Open("arcalive.db")
	if err != nil {
		log.Fatal("DB open:", err)
	}
	defer db.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	r := router.New(db)
	log.Println("arcalive server on :" + port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

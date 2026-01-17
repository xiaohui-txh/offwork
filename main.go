package main

import (
	"log"

	"offwork-backend/db"
	"offwork-backend/router"
)

func main() {
	log.Println("begin run")
	if err := db.InitDB(); err != nil {
		log.Fatalf("init db failed: %v", err)
	}

	r := router.InitRouter()
	r.Run(":8080")
}

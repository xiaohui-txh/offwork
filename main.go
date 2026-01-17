package main

import (
	"log"

	"offwork-backend/db"
	"offwork-backend/router"
)

func main() {
	log.Println("begin main")
	if err := db.InitDB(); err != nil {
		log.Fatalf("init db failed: %v", err)
	}

	log.Println("begin run...")
	r := router.InitRouter()
	r.Run(":8080")
}

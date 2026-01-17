package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() error {
	dsn := "root:f9rMNjuz@tcp(127.0.0.1:3306)/offwork?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("open mysql failed for %v", err)
		return err
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.Ping(); err != nil {
		log.Printf("ping mysql failed for %v", err)
		return err
	}

	log.Println("init mysql success")
	DB = db
	return nil
}

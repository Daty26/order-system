package db

import (
	"database/sql"
	"github.com/Daty26/order-system/order-service/internal/config"
	"log"
)

var DataB *sql.DB

func InitDB() {
	var err error
	DataB, err = sql.Open("postgres", config.GetDBConnectionString())
	if err != nil {
		log.Fatalf("failed to conn to db: %v", err)
	}
	if err = DataB.Ping(); err != nil {
		log.Fatalf("Couldn;t ping db: %v", err)
	}
	log.Println("Successfully connected to db")
}

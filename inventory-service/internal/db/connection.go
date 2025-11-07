package db

import (
	"database/sql"
	"github.com/Daty26/order-system/inventory-service/internal/config"
	_ "github.com/lib/pq"
	"log"
)

var DataDB *sql.DB

func InitDB() {
	var err error
	DataDB, err = sql.Open("postgres", config.GetConnString())
	if err != nil {
		log.Fatalf("Couldn't connect to db: %v", err.Error())
	}
	if err = DataDB.Ping(); err != nil {
		log.Fatalf("Couldn't ping db: %v", err.Error())
	}
	log.Println("connected to db!")
}

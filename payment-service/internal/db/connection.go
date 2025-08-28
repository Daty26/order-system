package db

import (
	"database/sql"
	"log"

	"github.com/Daty26/order-system/payment-service/internal/config"
	_ "github.com/lib/pq"
)

var DataDB *sql.DB

func InitDB() {
	var err error
	DataDB, err = sql.Open("postgres", config.GetDBConnectionString())
	if err != nil {
		log.Fatalf("couldn't connect to db: %v ", err.Error())
	}
	if err = DataDB.Ping(); err != nil {
		log.Fatalf("couldn't ping db: %v", err.Error())
	}
	log.Println("Connection succeed")
}

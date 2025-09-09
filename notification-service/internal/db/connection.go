package db

import (
	"database/sql"
	"fmt"
	"github.com/Daty26/order-system/notification-service/internal/config"
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
		log.Fatalf("Couldn';t ping the db: %v", err)
	}
	fmt.Println("Connected!")
}

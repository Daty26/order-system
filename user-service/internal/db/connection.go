package db

import (
	"database/sql"
	"github.com/Daty26/order-system/user-service/internal/config"
	_ "github.com/lib/pq"
	"log"
)

var DBConn *sql.DB

func DBInit() {
	var err error
	DBConn, err = sql.Open("postgres", config.GetConnString())
	if err != nil {
		log.Fatalf("Couldn't connect tp db: %v", err.Error())
	}
	if err = DBConn.Ping(); err != nil {
		log.Fatalf("Couldn't ping the db: %v \n", err.Error())
	}
}

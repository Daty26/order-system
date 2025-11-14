package config

import "fmt"

var (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1234"
	dbName   = "users"
)

func GetConnString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
}

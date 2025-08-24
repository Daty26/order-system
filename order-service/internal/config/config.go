package config

import "fmt"

var (
	Host     = "localhost"
	Port     = 5432
	User     = "postgres"
	Password = "1234"
	DBName   = "orders"
)

func GetDBConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", Host, Port, User, Password, DBName)
}

package config

import "fmt"

// chnage to vcar
const (
	Host     = "localhost"
	port     = 5432
	User     = "postgres"
	Password = "1234"
	DBName   = "payments"
)

func GetDBConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", Host, port, User, Password, DBName)
}

package config

import "fmt"

// chnagew to env var
const (
	Host     = "localhost"
	Port     = 5432
	User     = "postgres"
	Password = "1234"
	DBName   = "notification"
)

func GetConnString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", Host, Port, User, Password, DBName)
}

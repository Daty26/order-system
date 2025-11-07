package config

import "fmt"

const (
	Host   = "localhost"
	Port   = 5432
	User   = "postgres"
	Pass   = "1234"
	DBName = "inventory"
)

func GetConnString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", Host, Port, User, Pass, DBName)
}

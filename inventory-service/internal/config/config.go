package config

import (
	"fmt"
	"os"
	"strconv"
)

var (
	host     = getEnv("DB_HOST", "localhost")
	port, _  = strconv.Atoi(getEnv("DB_PORT", "5434"))
	user     = getEnv("DB_USER", "postgres")
	password = getEnv("DB_PASSWORD", "postgres")
	dbName   = getEnv("DB_NAME", "inventory")
)

func GetConnString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
}

func getEnv(value string, defaultVal string) string {
	if val := os.Getenv(value); val != "" {
		return val
	}
	return defaultVal
}

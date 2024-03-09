package config

import (
	"fmt"
	"os"
)

const LogLevel = "Debug" /* Debug | Info */

func GetConnectionString() string {
	dbname := getEnv("POSTGRES_DB", "empl")
	host := getEnv("POSTGRES_HOST", "localhost")
	port := getEnv("POSTGRES_PORT", "5433")
	user := getEnv("POSTGRES_USER", "user")
	passw := getEnv("POSTGRES_PASSWORD", "1234")
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, passw, dbname)
}

func GetMongoString() string {
	return "mongodb://root:example@localhost:27017/"
}

func getEnv(name string, defaultValue string) string {
	val, ok := os.LookupEnv(name)
	if ok {
		return val
	}
	return defaultValue
}

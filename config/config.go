package config

import (
	"fmt"
	"os"
	"time"
)

const LogLevel = "Debug" /* Debug | Info */

type PgConnectionOpt struct {
	Dbname string
	Host   string
	Port   string
	User   string
	Passw  string
}

func GetConnectionOpt() PgConnectionOpt {
	opt := PgConnectionOpt{}
	opt.Dbname = getEnv("POSTGRES_DB", "Empl")
	opt.Host = getEnv("POSTGRES_HOST", "localhost")
	opt.Port = getEnv("POSTGRES_PORT", "5433")
	opt.User = getEnv("POSTGRES_USER", "user")
	opt.Passw = getEnv("POSTGRES_PASSWORD", "1234")
	return opt
}

func GetConnectionString() string {
	opt := GetConnectionOpt()
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		opt.Host, opt.Port, opt.User, opt.Passw, opt.Dbname)
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

func GetRedisAddr() string {
	return getEnv("REDIS_ADDR", "localhost:6379")
}

const (
	RedisEmployeeDb         = 0
	RedisPassword           = ""
	RedisExpiration         = time.Minute * 1
	RedisResponseExpiration = time.Second * 30

	AppAddr = ":8080"
)

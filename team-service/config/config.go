package config

import (
	"os"
)

type Config struct {
	Port               string
	Host               string
	DBHost, DBPort     string
	DBUser, DBPass     string
	DBName             string
	AccessTokenSecret  string
	RefreshTokenSecret string
}

func LoadConfig() *Config {
	return &Config{
		Port:               os.Getenv("PORT"),
		Host:               os.Getenv("HOST"),
		DBHost:             os.Getenv("DB_HOST"),
		DBPort:             os.Getenv("DB_PORT"),
		DBUser:             os.Getenv("DB_USER"),
		DBPass:             os.Getenv("DB_PASS"),
		DBName:             os.Getenv("DB_NAME"),
		AccessTokenSecret:  os.Getenv("ACCESS_TOKEN_SECRET"),
		RefreshTokenSecret: os.Getenv("REFRESH_TOKEN_SECRET"),
	}
}

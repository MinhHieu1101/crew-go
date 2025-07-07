package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config holds all env variables
type Config struct {
	ENV                string
	Host               string
	DBHost             string
	DBPort             int
	DBUser             string
	DBPass             string
	DBName             string
	AccessTokenSecret  string
	RefreshTokenSecret string
	UserServicePort    string
	TeamServicePort    string
}

var Cfg *Config

func Init() {
	viper.SetConfigFile("../.env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No .env file found: %v", err)
	}

	Cfg = &Config{
		ENV:                viper.GetString("GO_ENV"),
		Host:               viper.GetString("HOST"),
		DBHost:             viper.GetString("DB_HOST"),
		DBPort:             viper.GetInt("DB_PORT"),
		DBUser:             viper.GetString("DB_USER"),
		DBPass:             viper.GetString("DB_PASS"),
		DBName:             viper.GetString("DB_NAME"),
		AccessTokenSecret:  viper.GetString("ACCESS_TOKEN_SECRET"),
		RefreshTokenSecret: viper.GetString("REFRESH_TOKEN_SECRET"),
		UserServicePort:    viper.GetString("USER_SERVICE_PORT"),
		TeamServicePort:    viper.GetString("TEAM_SERVICE_PORT"),
	}
}

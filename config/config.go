package config

import (
	env "github.com/Netflix/go-env"
	"log"
)

type Config struct {
	JwtAccessSecret     string `env:"JWT_ACCESS_SECRET"`
	JwtRefreshSecret    string `env:"JWT_REFRESH_SECRET"`
	AccessExpiryMinutes int    `env:"ACCESS_EXPIRY_MINUTES"`
	RefreshExpiryHours  int    `env:"REFRESH_EXPIRY_HOURS"`
	DbName string `env:"DB_DBNAME"`
	DbHost string `env:"DB_HOST"`
	DbUser string `env:"DB_USER"`
	DbPassword string `env:"DB_PASSWORD"`
	DbPort string `env:"DB_PORT"`
}

func GetConfig() Config {
	var cfg Config
	_, err := env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

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
}

func GetConfig() Config {
	var cfg Config
	_, err := env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

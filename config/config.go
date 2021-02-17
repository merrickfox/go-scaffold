package config

import (
	env "github.com/Netflix/go-env"
	"log"
)

type Config struct {
	JwtAccessSecret string `env:"JWT_ACCESS_SECRET"`
	JwtRefreshSecret string `env:"JWT_REFRESH_SECRET"`
}

func GetConfig() Config {
	var cfg Config
	_, err := env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}
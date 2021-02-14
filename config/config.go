package config

import (
	env "github.com/Netflix/go-env"
	"log"
)

type Config struct {
	JwtSecret string `env:"JWT_SECRET"`
}

func GetConfig() Config {
	var cfg Config
	_, err := env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}
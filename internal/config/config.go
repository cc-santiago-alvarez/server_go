package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type cfg struct {
	Env *ApiEnv
}

func InitConfig() *cfg {
	return &cfg{}
}

func (c *cfg) LoadEnvs() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	envFile := fmt.Sprintf(".env.%s", os.Getenv("ENV"))
	if err := godotenv.Load(envFile); err != nil {
		log.Fatal("Error loading env file: ", envFile)
	}

	log.Printf("Env loaded: %s\n", envFile)

	c.Env = &ApiEnv{
		port: os.Getenv("PORT"),
	}
}

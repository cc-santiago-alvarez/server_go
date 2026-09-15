package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env *ApiEnv
}

func InitConfig() *Config {
	return &Config{}
}

// Carga
func (c *Config) LoadEnvs() {
	if env := os.Getenv("ENV"); env != "" {
		envFile := fmt.Sprintf(".env.%s", env)
		if err := godotenv.Load(envFile); err != nil {
			log.Fatalf("Error loading env file %s: %v", envFile, err)
		}
		log.Printf("Env loaded: %s", envFile)
	}

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	c.Env = &ApiEnv{
		Port:     os.Getenv("PORT"),
		URI:      os.Getenv("MONGO_URI"),
		Database: os.Getenv("MONGO_DATABASE"),
	}
}

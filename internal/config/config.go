package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv(path string) {
	err := godotenv.Load(path)
	if err != nil {
		log.Fatalf("Error loading env file %s: %v", path, err)
	}
}

func QdrantHost() string {
	host := os.Getenv("QDRANT_HOST")
	if host == "" {
		host = "http://localhost:6333"
	}
	return host
}

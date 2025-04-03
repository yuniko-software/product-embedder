package main

import (
	"fmt"
	"log"
	"os"

	"product-embedder/internal/api"
	"product-embedder/internal/config"
	"product-embedder/internal/qdrant"
	"product-embedder/internal/utils"
)

func main() {
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = "dev.env"
	}
	config.LoadEnv(envFile)

	if err := qdrant.CreateCollection(); err != nil {
		log.Fatalf("Failed to create Qdrant collection: %v", err)
	}

	products, err := utils.LoadProductsCSV("data/products.csv")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Loaded %d products\n", len(products))

	qdrant.InsertAllProducts(products)
	api.SetupAPI()
}

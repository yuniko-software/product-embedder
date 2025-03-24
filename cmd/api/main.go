package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"

	"product-embedder/internal/config"
	"product-embedder/internal/models"
	"product-embedder/internal/qdrant"
	"product-embedder/internal/utils"
)

func main() {
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = "dev.env" // fallback
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

	var MaxWorkers = runtime.NumCPU() * 2
	wg := sync.WaitGroup{}
	jobs := make(chan models.Product, MaxWorkers)

	for i := 0; i < MaxWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for p := range jobs {
				input := p.ToEmbeddingInput()
				embedding, err := qdrant.GetEmbedding(input)
				if err != nil {
					log.Printf("[Worker %d] Failed to embed product %s: %v", workerID, p.ID, err)
					continue
				}
				err = qdrant.InsertProduct(p.ID, embedding, p)
				if err != nil {
					log.Printf("[Worker %d] Failed to insert %s: %v", workerID, p.ID, err)
				} else {
					log.Printf("[Worker %d] Inserted product: %s", workerID, p.ID)
				}
			}
		}(i + 1)
	}

	for _, p := range products {
		jobs <- p
	}
	close(jobs)
	wg.Wait()
	fmt.Println("✅ All products processed.")
}

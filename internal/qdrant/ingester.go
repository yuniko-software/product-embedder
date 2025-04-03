package qdrant

import (
	"fmt"
	"log"
	"runtime"
	"sync"

	"product-embedder/internal/models"
)

func InsertAllProducts(products []models.Product) {
	fmt.Printf("Inserting %d products with parallel processing...\n", len(products))

	workerCount := runtime.NumCPU() * 2
	jobs := make(chan models.Product, workerCount)
	wg := sync.WaitGroup{}

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for p := range jobs {
				input := p.ToEmbeddingInput()
				embedding, err := GetEmbedding(input)
				if err != nil {
					log.Printf("[Worker %d] Failed to embed product %s: %v", workerID, p.ID, err)
					continue
				}
				err = InsertProduct(p.ID, embedding, p)
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

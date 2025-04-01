package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"

	"product-embedder/internal/config"
	"product-embedder/internal/models"
	"product-embedder/internal/qdrant"
	"product-embedder/internal/utils"

	"github.com/gofiber/fiber/v2"
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

	setupAPI()
}

func setupAPI() {
	app := fiber.New()

	app.Get("/search", func(c *fiber.Ctx) error {
		query := c.Query("q")
		if query == "" {
			return c.Status(400).SendString("Missing query param `q`")
		}

		topK := 5
		if top := c.Query("top"); top != "" {
			if parsed, err := strconv.Atoi(top); err == nil && parsed > 0 {
				topK = parsed
			}
		}

		var maxPrice *float64
		if price := c.Query("maxPrice"); price != "" {
			if p, err := strconv.ParseFloat(price, 64); err == nil {
				maxPrice = &p
			}
		}

		results, err := qdrant.SearchProducts(query, topK, maxPrice)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.JSON(results)
	})

	log.Fatal(app.Listen(":8080"))
}

package qdrant

import (
	"fmt"
	"product-embedder/internal/models"

	"github.com/go-resty/resty/v2"
)

func CreateCollection() error {
	client := resty.New()
	body := map[string]interface{}{
		"vectors": map[string]interface{}{
			"size":     1536,
			"distance": "Cosine",
		},
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Put("http://localhost:6333/collections/products")

	if err != nil {
		return err
	}

	fmt.Println("Collection creation response:", resp.Status())
	return nil
}

func InsertProduct(id string, embedding []float32, p models.Product) error {
	client := resty.New()

	body := map[string]interface{}{
		"points": []map[string]interface{}{
			{
				"id":     id,
				"vector": embedding,
				"payload": map[string]interface{}{
					"name":           p.Name,
					"description":    p.Description,
					"price":          p.Price,
					"price_currency": p.PriceCurrency,
					"supply_ability": p.SupplyAbility,
					"minimum_order":  p.MinimumOrder,
				},
			},
		},
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Put("http://localhost:6333/collections/products/points")

	if err != nil {
		return err
	}

	fmt.Printf("Inserted %s: %s\n", id, resp.Status())
	return nil
}

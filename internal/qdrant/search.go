package qdrant

import (
	"encoding/json"
	"product-embedder/internal/config"

	"github.com/go-resty/resty/v2"
)

type SearchResult struct {
	ID      interface{}    `json:"id"`
	Payload map[string]any `json:"payload"`
	Score   float64        `json:"score"`
}

func SearchProducts(query string, topK int, maxPrice *float64) ([]SearchResult, error) {
	embedding, err := GetEmbedding(query)
	if err != nil {
		return nil, err
	}

	request := map[string]interface{}{
		"vector":       embedding,
		"top":          topK,
		"with_payload": true,
		"with_vector":  false,
	}

	if maxPrice != nil {
		request["filter"] = map[string]interface{}{
			"must": []map[string]interface{}{
				{
					"key": "price",
					"range": map[string]interface{}{
						"lt": *maxPrice,
					},
				},
			},
		}
	}

	host := config.QdrantHost()

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		Post(host + "/collections/products/points/search")

	if err != nil {
		return nil, err
	}

	var result struct {
		Result []SearchResult `json:"result"`
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}

	return result.Result, nil
}

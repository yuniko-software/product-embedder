$env:ENV_FILE = "prod.env"
docker-compose up --build / docker-compose up -d qdrant
go run ./cmd/api
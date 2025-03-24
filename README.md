$env:ENV_FILE = "prod.env"
docker-compose up --build
go run ./cmd/api
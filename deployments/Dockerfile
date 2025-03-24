FROM golang:1.21-bookworm

WORKDIR /app

ENV GOTOOLCHAIN=auto

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o product-embedder ./cmd/api

CMD ["./product-embedder"]


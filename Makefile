# Start Docker infrastructure
infra:
	docker-compose up -d

# Run the Ingestion API
api:
	go run cmd/producer/main.go

# Run the Consumer Worker
worker:
	go run cmd/worker/main.go

# Download and clean up dependencies
tidy:
	go mod tidy

# Stop all containers
clean:
	docker-compose down
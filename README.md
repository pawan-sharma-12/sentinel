# Sentinel: Distributed Event Alerting & Rate-Limiting Engine

Sentinel is a high-performance, distributed alerting system built in **Go**. It is designed to ingest massive volumes of events via a REST API, buffer them using **Apache Kafka**, and process them through workers that implement distributed rate-limiting using **Redis**.

## 🚀 Key Features
- **Asynchronous Processing:** Decouples event ingestion from notification delivery using Kafka.
- **Distributed Rate Limiting:** Implements a Fixed Window algorithm using Redis to prevent notification spam across multiple worker instances.
- **Graceful Shutdown:** Handles OS signals (SIGINT/SIGTERM) to ensure no data is lost during worker deployment or scaling.
- **Clean Architecture:** Separated into `cmd` (entry points), `internal/cache` (Redis), and `internal/models`.

## 🛠 Tech Stack
- **Language:** Go (Golang)
- **API Framework:** Gin Gonic
- **Message Broker:** Apache Kafka
- **State Store:** Redis
- **Infrastructure:** Docker & Docker Compose

## 🏗 System Architecture
1. **Producer (API):** Receives JSON payloads and pushes them into a Kafka topic (`raw-events`).
2. **Message Bus (Kafka):** Acts as a shock absorber for traffic spikes.
3. **Worker (Consumer):** Pulls messages from Kafka.
4. **Rate Limiter (Redis):** Tracks user activity. If a user exceeds 5 requests per minute, the worker drops the event.
5. **Notification:** If within limits, the worker simulates sending an alert (Email/SMS/Push).



## 🚥 Getting Started

### Prerequisites
- Docker & Docker Desktop
- Go 1.21+

### 1. Setup Infrastructure
```bash
make infra
**2. Run the Services**
Open two terminals:

Bash

# Terminal 1: Start the Worker
go run cmd/worker/main.go

# Terminal 2: Start the API
go run cmd/producer/main.go
**3. Test the Rate Limiter**
Send 6 requests for the same user to see the 6th one get throttled:

Bash

for i in {1..6}; do
  curl -X POST http://localhost:8080/v1/alerts \
  -H "Content-Type: application/json" \
  -d '{"user_id": "pawan_sharma", "type": "PRICE_DROP", "payload": {"item": "MacBook"}}'
done

📈 Future Improvements
[ ] Implement Sliding Window rate limiting using Redis Sorted Sets.

[ ] Add Protobuf support for more efficient Kafka message serialization.

[ ] Implement Manual Offsets in Kafka for exactly-once processing guarantees.

Author: Pawan Sharma

Role: Technical Lead @ HCL Technologies

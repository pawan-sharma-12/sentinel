Markdown

# Sentinel: Distributed Event Alerting & Rate-Limiting Engine

Sentinel is a high-performance, distributed alerting system built in **Go**. It is designed to ingest massive volumes of events (1billion+ req/day) via a REST API, buffer them using **Apache Kafka**, and process them through workers that implement distributed rate-limiting using **Redis**.

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
2. Run the Services
Open two terminals:

Bash

# Terminal 1: Start the Worker
go run cmd/worker/main.go

# Terminal 2: Start the API
go run cmd/producer/main.go
3. Test the Rate Limiter
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

Location: Bengaluru, India


### Stress Test :

hey -n 1000 -c 20 -m POST \
-H "Content-Type: application/json" \
-d '{"user_id": "load_test_user", "type": "STRESS_TEST", "payload": {"data": "peak_test"}}' \
http://localhost:8080/v1/alerts

Summary:
  Total:        0.0728 secs
  Slowest:      0.0266 secs
  Fastest:      0.0001 secs
  Average:      0.0014 secs
  Requests/sec: 13735.5798
  
  Total data:   120000 bytes
  Size/request: 120 bytes

Response time histogram:
  0.000 [1]     |
  0.003 [948]   |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.005 [11]    |
  0.008 [20]    |■
  0.011 [0]     |
  0.013 [0]     |
  0.016 [0]     |
  0.019 [0]     |
  0.021 [0]     |
  0.024 [1]     |
  0.027 [19]    |■


Latency distribution:
  10% in 0.0003 secs
  25% in 0.0004 secs
  50% in 0.0007 secs
  75% in 0.0012 secs
  90% in 0.0019 secs
  95% in 0.0029 secs
  99% in 0.0247 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0001 secs, 0.0001 secs, 0.0266 secs
  DNS-lookup:   0.0001 secs, 0.0000 secs, 0.0029 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0020 secs
  resp wait:    0.0012 secs, 0.0001 secs, 0.0174 secs
  resp read:    0.0000 secs, 0.0000 secs, 0.0008 secs

Status code distribution:
  [202] 1000 responses


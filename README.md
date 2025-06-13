🧠 Objective:
Build a high-performance system in Go to simulate load over TCP using Redis and controlled concurrency.

✅ Requirements Breakdown:
🧩 1. Redis + TPS Loader
Read TPS config (e.g., 100 records/sec).

Fetch records from Redis using keys like record:1, record:2, etc.

Send them at the defined rate every second.

🔌 2. TCP Client
Open up to 10 TCP connections to the server.

Each connection sends max 100 requests (records).

Waits for response → maps it back to the original Redis record.

🎧 3. TCP Server
Listens on 10 different ports.

Processes incoming data and sends back a unique response ID for each request.

🧵 4. Concurrency
Use Goroutines, Channels, and Worker Pools to:

Control rate (TPS)

Manage concurrent sending

Limit requests per connection

Handle responses efficiently

⚙️ Goal:
Achieve controlled load generation using concurrent TCP sockets, ensuring rate-limited processing and proper mapping of responses back into Redis.


# 🚀 Go Redis Record Generator

A high-performance Golang application to generate and store up to 500,000 records in Redis using a REST API with Gin framework. It uses Redis pipelines to speed up data insertion and also provides an endpoint to count total records.

---

## 📦 Features

- Generate a custom number of records via HTTP POST
- Uses Redis pipeline for fast batch writes
- Count total records stored using SCAN
- Minimal setup, clean and modular codebase

---

## 🛠️ Requirements

- Go 1.18 or above
- Redis Server (running locally or remotely)
- Modules installed via:

```bash
go mod tidy

🚀 Running the Server
Start your Redis server locally (default localhost:6379)

Run the app:go run main.go

🔄 API Endpoints
📤 Generate Records
Generates the specified number of Redis records.

URL:
POST /generate
{
  "count": 500000
}
Response
{
  "message": "500000 records stored successfully",
  "duration": "5.3212342s"
}

📊 Count Records
Counts how many records exist in Redis under the pattern record:*.

URL:
GET /count

Response:

{
  "total_records": 500000
}


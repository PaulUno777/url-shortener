# 📦 URL Shortener API

A simple and efficient URL shortener built with **Go**, **Fiber**, and **Redis**.  
Supports custom short URLs, expiration times, and basic rate limiting per IP.

---

## 🚀 Features

- Shorten any valid URL
- Set custom aliases (optional)
- Expiration support for shortened URLs
- Rate limiting per IP using Redis
- Clean REST API (JSON-based)
- Dockerized environment for easy setup

---

## 🛠️ Tech Stack

- [Go](https://golang.org/)
- [Fiber](https://gofiber.io/) - Web framework
- [Redis](https://redis.io/) - In-memory DB for rate limiting and URL mapping
- [Docker](https://www.docker.com/) - Containerized setup

---

## 📂 Project Structure

.
├── api
│ ├── main.go
│ ├── routes/
│ └── ...
├── db
│ └── Dockerfile
├── .env.example
├── docker-compose.yml
└── README.md

## 🧪 Endpoints

### ➕ Shorten URL

**POST** `/api/v1`  
Creates a shortened version of the input URL.

**Request JSON:**

```json
{
  "url": "https://example.com",
  "short": "custom123", // optional
  "expiry": 3600 // optional (in seconds)
}
```

**Response JSON:**

```json
{
  "url": "https://example.com",
  "short": "/custom123",
  "expiry": 3600,
  "rate_limit": 9,
  "rate_limit_rest": 29
}
```

### 🔗 Resolve Short URL

**GET** `/:url`  
Redirects to the original long URL.

Example:
GET /custom123 → Redirects to https://example.com

## ⚙️ Environment Variables

Create a .env file (based on `.env.example`)

## 🐳 Docker Setup

```bash
$ docker compose up -d
```

The API will be accessible at: http://localhost:3000

##  📌 Notes
- Rate limiting is per IP with a 30-minute window.
- Default expiry for shortened URLs is 24 hours unless specified.
- Custom short aliases must be unique or the API returns a conflict.

## 📃 License
MIT © 2025 PaulUno777
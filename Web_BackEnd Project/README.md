# 🌐 Web Backend Project — Node.js Web App + Go Microservices (Experimental)

This repository contains a **hybrid backend project** with two separate implementations:

1) **Node.js + Express + EJS** — main web application (pages, auth, CRUD, API proxy endpoints).  
2) **Go (HTTP + gRPC + NATS + Redis)** — experimental microservice-style backend (Clean Architecture structure).

> ✅ The Node.js application is the primary runnable part.  
> ⚠️ The Go backend is an experimental distributed-systems implementation and is **not fully integrated** with the Node app yet.

---

## 🧩 Project Structure

### Node.js (Web App)
- `app.js` — main Express server, sessions, MongoDB connection, routes
- `routes/` — route definitions (auth, main, landmarks, quiz, API proxy, transactions)
- `controllers/` — business logic (auth, CRUD, external API calls)
- `models/` — Mongoose schemas
- `views/` — EJS templates (UI pages)
- `public/` — static files (CSS/JS/uploads)

### Go (Distributed Backend — Experimental)
- `main.go` — HTTP server + gRPC startup
- `internal/` — Clean Architecture layers (domain / usecase / repository / delivery)
- `proto/`, `grpc/` — gRPC contracts and server implementations
- `transaction/`, `quiz/` — service modules (microservice-like structure)

---

## ✨ Key Features (Implemented)

### ✅ Node.js Web App
- Authentication + sessions
- CRUD for main entities (e.g., landmarks / users / transactions)
- Server-side rendering with EJS
- API proxy endpoints to external services (news / stocks / etc.)
- Quiz page and result storage (MongoDB)

### ✅ Go Backend (Systems / Microservices Focus)
- gRPC server (user CRUD style services)
- NATS messaging integration
- Redis caching layer (optional)
- Multi-service style ports (HTTP + gRPC + service endpoints)
- Clean Architecture structure for maintainability

---

## 🔒 Security Note (Important)

This project uses environment variables for credentials and API keys.

✅ The repository includes `env.example` (template).  
❌ Do NOT commit real secrets (`.env` is ignored).

If any credentials were previously committed, they should be rotated immediately.

---

## ▶️ Run (Node.js — Main App)

 - npm install
 - npm run dev

## Open:

 - http://localhost:3000

## 🧪 Go Backend (Optional / Experimental)

 - Requires MongoDB and NATS running locally.
 - This part is currently used as an architecture experiment.

 - go mod download
 - go run main.go


## Expected ports (may vary by config):

- HTTP: 8080

- gRPC: 50051

# other services: 8081, 8082

## 📌 Future Improvements

Integrate Node gateway with Go services via gRPC

Remove hardcoded config and use .env everywhere

Add Docker Compose for MongoDB + NATS + Redis

Add tests (unit/integration)

CI/CD pipeline (GitHub Actions)

Centralized structured logging + metrics

## 📄 License

Educational / Demonstration project.

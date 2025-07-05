# Subscription Service
A concurrent Subscription Service developed in Go. This project demonstrates extensive use of concurrency patterns. It runs email sending, PDF generation, and invoice creation in parallel, synchronizes tasks with wait groups and communicates between goroutines using channels.

---

## Features

- **Choose a subscription plan**
- **Send activation emails to users**
- **Generate invoices upon subscription**
- **Generate PDF manuals for subscriptions**
- **All major processes implemented using Go's concurrency primitives**
- **RESTful API powered by Chi router**
- **Session management with Redis**
- **Database persistence with PostgreSQL and pgx driver**
- **Containerized with Docker for easy setup**
- **Comprehensive unit tests for handlers, routes, and concurrency logic**

---

## Used Tools:

- [Golang](https://github.com/golang/go)
- [Docker Compose](https://github.com/docker/compose)
- [Chi Router](https://github.com/go-chi/chi)
- [PGX PostgreSQL driver](https://github.com/jackc/pgx)
- [Redis](https://github.com/redis/redis)
- [MailHog](https://github.com/mailhog/MailHog)

---

## How to run?

**Run all this commands in the root level of project**

### 1- Install Go Dependencies

```bash
go mod download
```

### 2- Run Database

```bash
docker-compose up -d
```

### 3- Execute db.sql File

### 4- Run The Back-end API

```bash
go run ./cmd/web/.
```
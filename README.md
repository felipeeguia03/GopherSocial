# GopherSocial

A RESTful social network API built with Go. Users can register, create posts, follow each other, comment, and get a personalized feed.

**Live API:** [gophersocial.onrender.com](https://gophersocial.onrender.com/v1/health)  
**Frontend:** [gopher-social-web.vercel.app](https://gopher-social-web.vercel.app)  
**Docs:** [gophersocial.onrender.com/v1/swagger](https://gophersocial.onrender.com/v1/swagger/index.html)

---

## Features

- JWT authentication with email-based registration flow
- Account confirmation via transactional email (Brevo)
- Posts with comments and personalized feed based on followed users
- Follow / unfollow system with suggested users
- Role-based access control (user / moderator / admin)
- Redis caching layer for user lookups
- Rate limiting middleware
- Swagger documentation
- CI via GitHub Actions (build, vet, staticcheck, tests)
- CD via Render (auto-deploy on push to `main`)

## Tech Stack

| Layer | Technology |
|-------|------------|
| Language | [Go](https://go.dev/) |
| Router | [chi](https://github.com/go-chi/chi) |
| Database | [PostgreSQL 16](https://www.postgresql.org/) |
| Cache | [Redis](https://redis.io/) via [go-redis](https://github.com/redis/go-redis) |
| Auth | [golang-jwt/jwt](https://github.com/golang-jwt/jwt) |
| Migrations | [golang-migrate](https://github.com/golang-migrate/migrate) |
| Email | [Brevo](https://www.brevo.com/) |
| Docs | [Swagger / swaggo](https://github.com/swaggo/swag) |
| Logging | [zap](https://github.com/uber-go/zap) |
| Validation | [go-playground/validator](https://github.com/go-playground/validator) |
| Containers | [Docker](https://www.docker.com/) |

## Project Structure

```
.
├── cmd/
│   ├── api/          # HTTP handlers, middleware, routing
│   └── migrate/      # Migrations and seed data
├── internal/
│   ├── auth/         # JWT authenticator
│   ├── db/           # Database connection
│   ├── env/          # Environment variable helpers
│   ├── mailer/       # Email client and templates
│   ├── ratelimiter/  # Rate limiting
│   └── store/        # Data access layer (PostgreSQL + Redis)
└── docs/             # Auto-generated Swagger docs
```

## Endpoints

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/v1/auth/register` | Register a new user | No |
| POST | `/v1/auth/token` | Login and get JWT | No |
| PUT | `/v1/users/activate/{token}` | Activate account via email | No |
| GET | `/v1/users/{id}` | Get user profile | JWT |
| GET | `/v1/users/search?q=` | Search users by username | JWT |
| GET | `/v1/users/suggested` | Get suggested users to follow | JWT |
| PUT | `/v1/users/{id}/follow` | Follow a user | JWT |
| PUT | `/v1/users/{id}/unfollow` | Unfollow a user | JWT |
| GET | `/v1/feed` | Get personalized feed | JWT |
| POST | `/v1/posts` | Create a post | JWT |
| GET | `/v1/posts/{id}` | Get a post | JWT |
| PATCH | `/v1/posts/{id}` | Update a post | JWT (owner/moderator) |
| DELETE | `/v1/posts/{id}` | Delete a post | JWT (owner/admin) |
| POST | `/v1/posts/{id}/comments` | Add a comment | JWT |

## Getting Started

### Prerequisites

- Go 1.21+
- Docker
- [golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

### Running locally

```bash
# Start PostgreSQL and Redis
docker compose up -d

# Copy and fill environment variables
cp .envrc.example .envrc
source .envrc

# Run migrations
make migrate-up

# Seed the database (optional)
make seed

# Start the server
go run ./api
```

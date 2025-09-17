# GoSocial

A personal project to learn Go by building a social media platform with JWT authentication, posts, and a Next.js frontend.

## Learning Goals

- Go backend development with Chi router
- JWT authentication and middleware
- PostgreSQL database operations

## Tech Stack

- **Backend**: Go, Chi Router, PostgreSQL, JWT
- **Frontend**: Next.js, TypeScript, Tailwind CSS
- **Database**: PostgreSQL with Docker

## Prerequisites

- Go 1.25+
- Node.js 20+
- Docker

## � How to Run

### 1. Start Database
```bash
docker-compose up -d
```

### 2. Initialize Database
```bash
psql -h localhost -U admin -d social -f scripts/db_init.sql
```

### 3. Run Backend
```bash
go run cmd/api/*.go
```
API runs on `http://localhost:8080`

### 4. Run Frontend
```bash
cd web
npm install
npm run dev
```
Frontend runs on `http://localhost:3000`

## Current API Endpoints

- `POST /v1/auth/register` - User registration
- `POST /v1/auth/login` - User login  
- `GET /v1/users/me` - Get user profile (auth required)
- `GET /v1/posts` - Get posts (auth required)
- `GET /v1/health` - Health check

## Planned Features

- [ ] Create/edit posts
- [ ] User profiles
- [ ] Comment system
- [ ] File uploads
- [ ] Real-time updates
- [ ] Search functionality

---

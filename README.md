# User Management Service

User Management Service.

## Tech Stack
- Go 1.26.3
- Gin Framework
- PostgreSQL 15
- JWT Authentication
- GORM

## Features (P0)
- ✅ User registration (admin only)
- ✅ User login with JWT
- ✅ User logout (whitelist)
- ✅ Get current user info
- ✅ CRUD users (admin only)
- ✅ Soft delete users
- ✅ Role-based access (admin, supervisor, warehouse_operator, viewer)
- ✅ Session whitelist (track active sessions)
- ✅ Role Management
- ✅ Menu Management
- ✅ Permission Management
- ✅ Login History Management
- ✅ Forgot Password
- ✅ Reset Password
- ✅ Password Strength Check
- ✅ Password Validation

## Prerequisites
- Go 1.26.3+
- PostgreSQL 15+
- Docker & Docker Compose (optional)

## Setup

### Local Development

1. Clone repository
2. Copy `.env.example` to `.env` and adjust values
3. Install Make ```choco install make```


### Run migrations

```bash
make migrate
```
### Run Seeder

```bash
make seed
```

### Generate Swagger

```bash
make swagger
```

### Run Server

```bash
make run
```
# Authentio

<div align="center">

**A modern, production-ready authentication API built for scale**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-ready-2496ED?style=flat&logo=docker)](https://www.docker.com/)

[Features](#features) • [Quick Start](#quick-start) • [API Reference](#api-reference) • [Architecture](#architecture)

</div>

---

## Overview

Authentio is a high-performance authentication and authorization API engineered with Go. It provides enterprise-grade security features including JWT-based authentication, OAuth2 social login, two-factor authentication, and distributed rate limiting—all packaged in a Docker-ready solution.

Perfect for microservices architectures, mobile apps, and modern web applications that need bulletproof auth without the overhead.

## Features

### Core Authentication
- **🔐 JWT-Based Auth** - Stateless access and refresh tokens with automatic rotation
- **🌐 OAuth2 Integration** - Google Sign-In support (extensible to other providers)
- **🔑 Password Management** - Secure reset flow with email-based verification
- **👤 User Management** - Complete CRUD operations for user profiles

### Security
- **🛡️ Two-Factor Authentication** - Email-based OTP for enhanced security
- **⚡ Rate Limiting** - Redis-powered distributed rate limiting
- **🚫 Token Blacklisting** - Instant token revocation support
- **🔒 Secure Defaults** - Bcrypt password hashing, HTTPS-ready

### Performance & Scalability
- **🚀 Go-Powered** - Concurrent request handling and minimal resource footprint
- **📦 Containerized** - Docker Compose setup for consistent environments
- **💾 PostgreSQL** - ACID-compliant data persistence
- **⚙️ Redis** - High-speed caching and session management

## Quick Start

### Prerequisites
- Docker & Docker Compose
- A Google Cloud project (for OAuth2)
- SMTP credentials (Gmail recommended for development)

### Installation

1. **Clone and navigate**
   ```bash
   git clone https://github.com/your-username/authentio.git
   cd authentio
   ```

2. **Configure environment**
   ```bash
   cp .env.example .env
   # Edit .env with your credentials (see configuration below)
   ```

3. **Launch the stack**
   ```bash
   docker-compose up --build
   ```

4. **Verify it's running**
   ```bash
   curl http://localhost:8080/health
   ```

The API will be accessible at `http://localhost:8080/api/v1`

## Configuration

Create a `.env` file with these variables:

```env
# Server
SERVER_PORT=8080
APP_ENV=development

# Database
POSTGRES_DSN=postgres://postgres:secret@db:5432/authentio_db?sslmode=disable

# Redis
REDIS_ADDR=redis:6379
REDIS_PASS=

# Security
JWT_SECRET=your-super-secret-jwt-key-min-32-chars

# Email (SMTP)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=noreply@yourdomain.com

# OAuth2 (Google)
GOOGLE_CLIENT_ID=your-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=GOCSPX-your-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/api/v1/auth/google/callback
```

**Security Note**: Use app-specific passwords for Gmail and never commit your `.env` file.

## Architecture

```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │
       ▼
┌─────────────┐      ┌──────────┐
│  Gin Router │◄────►│  Redis   │ (Rate Limiting, Blacklist)
└──────┬──────┘      └──────────┘
       │
       ▼
┌─────────────┐      ┌──────────┐
│  Handlers   │◄────►│PostgreSQL│ (User Data, Tokens)
└─────────────┘      └──────────┘
```

**Tech Stack:**
- **Framework**: Gin (fastest Go HTTP router)
- **Database**: PostgreSQL 15+ with GORM
- **Cache**: Redis 7+ for distributed operations
- **Authentication**: JWT with RS256 signing
- **Deployment**: Docker + Docker Compose

## API Reference

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication Endpoints

#### Register User
```http
POST /auth/register
Content-Type: application/json

{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john@example.com",
  "password": "SecurePass123!"
}
```

**Response (201)**
```json
{
  "user": {
    "id": 1,
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "is_active": true,
    "created_at": "2025-11-05T10:00:00Z"
  },
  "message": "Registration successful"
}
```

#### Login
```http
POST /auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "SecurePass123!"
}
```

**Response (200)**
```json
{
  "user": {
    "id": 1,
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "is_active": true
  },
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "a1b2c3d4e5f6...",
  "expires_in": 3600
}
```

#### Refresh Token
```http
POST /auth/refresh
Content-Type: application/json

{
  "refresh_token": "a1b2c3d4e5f6..."
}
```

#### Password Reset Flow
```http
# Step 1: Request reset
POST /auth/forgot-password
{
  "email": "john@example.com"
}

# Step 2: Reset with code
POST /auth/reset-password
{
  "email": "john@example.com",
  "code": "123456",
  "new_password": "NewSecurePass456!"
}
```

### OAuth2 Endpoints

#### Google Login (Frontend Flow)
```http
POST /auth/google/login
Content-Type: application/json

{
  "id_token": "google_id_token_from_gsi"
}
```

#### Google Login (Server-Side Flow)
```http
# Step 1: Redirect to Google
GET /auth/google/redirect

# Step 2: Handle callback (automatic)
GET /auth/google/callback?code=...
```

### Two-Factor Authentication

#### Enable 2FA
```http
POST /2fa/enableOtp
Authorization: Bearer {access_token}
```

#### Send OTP Code
```http
POST /2fa/sendOtp
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "email": "john@example.com"
}
```

#### Verify OTP
```http
POST /auth/2fa/verify
Content-Type: application/json

{
  "email": "john@example.com",
  "code": "123456"
}
```

#### Disable 2FA
```http
POST /2fa/disableOtp
Authorization: Bearer {access_token}
```

### User Management

#### Get Profile
```http
GET /user/getProfile
Authorization: Bearer {access_token}
```

#### Update Profile
```http
PUT /user/updateProfile
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "first_name": "Jonathan",
  "last_name": "Doer",
  "email": "jonathan@example.com"
}
```

### Error Responses

All errors follow this structure:
```json
{
  "error": "Error message",
  "code": "ERROR_CODE"
}
```

Common status codes:
- `400` - Bad Request (validation errors)
- `401` - Unauthorized (invalid/missing token)
- `404` - Not Found
- `409` - Conflict (duplicate email)
- `429` - Too Many Requests (rate limited)
- `500` - Internal Server Error

## Development

### Running Tests
```bash
go test ./... -v
```

### Database Migrations
```bash
# Migrations run automatically on startup
# To create new migrations, add to internal/database/migrations/
```

### Local Development Without Docker
```bash
# Ensure PostgreSQL and Redis are running locally
go mod download
go run cmd/api/main.go
```

## Security Best Practices

1. **Environment Variables**: Never commit `.env` files
2. **JWT Secret**: Use a strong, randomly generated secret (min 32 chars)
3. **HTTPS**: Always use HTTPS in production
4. **Rate Limiting**: Configure appropriate limits for your use case
5. **Password Policy**: Enforce strong passwords (8+ chars, mixed case, numbers, symbols)
6. **2FA**: Enable for sensitive operations
7. **Token Rotation**: Refresh tokens expire and should be rotated regularly

## Troubleshooting

### Database Connection Issues
```bash
# Check if PostgreSQL is running
docker-compose ps

# View database logs
docker-compose logs db
```

### Redis Connection Issues
```bash
# Verify Redis is accessible
docker-compose exec redis redis-cli ping
# Expected: PONG
```

### SMTP/Email Issues
- Verify SMTP credentials are correct
- For Gmail: Enable 2FA and create an App Password
- Check firewall rules for outbound port 587

### OAuth2 Issues
- Verify redirect URL matches Google Cloud Console exactly
- Ensure client ID and secret are correct
- Check that OAuth consent screen is configured

## Roadmap

- [ ] Multi-provider OAuth2 (GitHub, Facebook, Apple)
- [ ] TOTP-based 2FA (Google Authenticator)
- [ ] Role-based access control (RBAC)
- [ ] Session management dashboard
- [ ] Audit logging
- [ ] API key authentication
- [ ] WebAuthn/Passkey support

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

**Adebakin Olujimi**  
[@olujimi_the_dev](https://x.com/olujimi_the_dev)

---

<div align="center">
Built with ❤️ using Go and modern cloud-native technologies
</div>
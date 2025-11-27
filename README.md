# Authentio - Enterprise Authentication API

A production-ready authentication API built with Go, featuring JWT auth, OAuth2, 2FA, and distributed rate limiting.

**Base URL**: `http://localhost:8080/api/v1` | **Swagger Docs**: `http://localhost:8080/swagger/index.html`

---

## Quick Start

### Prerequisites

- Docker & Docker Compose
- Google Cloud credentials (for OAuth2)
- SMTP credentials (for email notifications)

### Setup

```bash
# Clone and setup
git clone <repo>
cd authentio

# Configure environment
cp .env.example .env
# Edit .env with your credentials

# Start services
docker-compose up -d

# Run application
go run cmd/server/main.go
```

Visit http://localhost:8080/swagger/index.html to test endpoints

---

## Project Structure

```
authentio/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
│
├── internal/
│   ├── config/
│   │   ├── config.go            # Configuration loader
│   │   └── google_oauth.go       # Google OAuth setup
│   │
│   ├── constants/
│   │   └── otp_types.go          # OTP type constants
│   │
│   ├── database/
│   │   ├── db.go                 # Database connection
│   │   ├── migrate.go            # Database migrations
│   │   ├── otp_repository.go     # OTP data access
│   │   ├── token_repository.go   # Token data access
│   │   ├── twofa_repository.go   # 2FA data access
│   │   └── user_repository.go    # User data access
│   │
│   ├── handler/
│   │   ├── auth_handler.go       # Auth endpoints
│   │   ├── handler.go            # Handler initialization
│   │   ├── requests.go           # Request DTOs
│   │   ├── twofa_handler.go      # 2FA endpoints
│   │   ├── user_handler.go       # User endpoints
│   │   └── validator.go          # Request validation
│   │
│   ├── middleware/
│   │   ├── auth_middleware.go    # JWT authentication
│   │   ├── blacklist_middleware.go # Token blacklisting
│   │   ├── cors.go               # CORS configuration
│   │   ├── logger.go             # Request logging
│   │   ├── ratelimit_inmem.go    # In-memory rate limiting
│   │   └── ratelimit_redis.go    # Redis rate limiting
│   │
│   ├── models/
│   │   ├── auth.go               # Auth models
│   │   ├── auth_response.go      # Auth responses
│   │   ├── model.go              # Base models
│   │   ├── otp.go                # OTP model
│   │   ├── refresh_token.go      # Refresh token model
│   │   ├── twoFA.go              # 2FA model
│   │   ├── user.go               # User model
│   │   └── user_profile.go       # User profile model
│   │
│   ├── repository/
│   │   ├── 2fa_repository.go     # 2FA repository interface
│   │   ├── otp_repository.go     # OTP repository interface
│   │   ├── token_repository.go   # Token repository interface
│   │   └── user_repository.go    # User repository interface
│   │
│   ├── router/
│   │   └── router.go             # Route definitions
│   │
│   └── service/
│       ├── auth_service.go       # Auth business logic
│       ├── token_service.go      # Token management
│       ├── twofa_service.go      # 2FA business logic
│       └── user_service.go       # User business logic
│
├── migrations/
│   ├── 001_init_schema.up.sql    # Initial schema
│   └── 001_init_schema.down.sql  # Schema rollback
│
├── pkg/
│   ├── email/
│   │   ├── email.go              # Email service interface
│   │   └── sendgrid.go           # SendGrid implementation
│   │
│   ├── jwt/
│   │   └── jwt.go                # JWT token management
│   │
│   ├── logger/
│   │   └── logger.go             # Structured logging
│   │
│   ├── password/
│   │   └── password.go           # Password hashing/verification
│   │
│   └── response/
│       └── response.go           # Response formatting
│
├── docs/
│   ├── swagger.json              # Swagger OpenAPI spec
│   ├── swagger.yaml              # Swagger YAML spec
│   └── docs.go                   # Swagger documentation
│
├── infra/
│   └── [Kubernetes manifests]    # K8s deployment files
│
├── .env                          # Environment variables (local)
├── docker-compose.yml            # Docker Compose setup
├── Dockerfile                    # Container image
├── go.mod                        # Go module definition
├── go.sum                        # Go dependencies
└── README.md                     # This file
```

### Key Directories

- **`cmd/`** - Application entry points
- **`internal/`** - Private application code (not importable by external packages)
  - `config/` - Configuration management
  - `database/` - Database layer (repositories, migrations)
  - `handler/` - HTTP request handlers
  - `middleware/` - HTTP middleware
  - `models/` - Data models
  - `router/` - Route definitions
  - `service/` - Business logic
- **`pkg/`** - Reusable packages that can be imported
- **`migrations/`** - Database schema migrations
- **`docs/`** - API documentation (Swagger/OpenAPI)
- **`infra/`** - Infrastructure as code (Docker, Kubernetes)

---

## API Endpoints

### Base URL

```
http://localhost:8080/api/v1
```

---

## Authentication Endpoints

### 1. Register User

**Request:**

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

**Success Response (201):**

```json
{
  "user": {
    "id": 1,
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "is_active": true,
    "created_at": "2025-11-27T10:00:00Z"
  },
  "message": "Registration successful"
}
```

**Error Response (400):**

```json
{
  "validation_error": {
    "email": "Invalid email format",
    "password": "Password must contain uppercase, lowercase, number, and special character"
  }
}
```

---

### 2. Login

**Request:**

```http
POST /auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "SecurePass123!"
}
```

**Success Response (200):**

```json
{
  "user": {
    "id": 1,
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com"
  },
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "a1b2c3d4e5f6g7h8i9j0...",
  "expires_in": 900
}
```

**Error Response (401):**

```json
{
  "error": "Invalid email or password"
}
```

---

### 3. Refresh Token

**Request:**

```http
POST /auth/refresh
Content-Type: application/json

{
  "refresh_token": "a1b2c3d4e5f6g7h8i9j0..."
}
```

**Success Response (200):**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "x9y8z7w6v5u4t3s2r1q0...",
  "expires_in": 900
}
```

---

### 4. Forgot Password

**Request:**

```http
POST /auth/forgot-password
Content-Type: application/json

{
  "email": "john@example.com"
}
```

**Success Response (200):**

```json
{
  "message": "Password reset code sent to your email"
}
```

---

### 5. Reset Password

**Request:**

```http
POST /auth/reset-password
Content-Type: application/json

{
  "email": "john@example.com",
  "code": "123456",
  "new_password": "NewSecurePass456!"
}
```

**Success Response (200):**

```json
{
  "message": "Password reset successfully"
}
```

---

## OAuth2 Endpoints

### 6. Google Login (Client-Side Flow)

**Request:**

```http
POST /auth/google/login
Content-Type: application/json

{
  "id_token": "eyJhbGciOiJSUzI1NiIsImtpZCI6IjEifQ..."
}
```

**Success Response (200):**

```json
{
  "user": {
    "id": 1,
    "email": "user@gmail.com",
    "first_name": "Google",
    "last_name": "User"
  },
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "a1b2c3d4e5f6g7h8i9j0...",
  "expires_in": 900
}
```

---

### 7. Google Redirect (Server-Side Flow)

**Request:**

```http
GET /auth/google/redirect
```

**Response:** Redirects to Google OAuth consent screen

---

### 8. Google Callback

**Request:**

```http
GET /auth/google/callback?code=4/0AX4XfWh...&state=state_value
```

**Response:** Redirects to frontend with tokens or error

---

## Two-Factor Authentication

### 9. Enable 2FA

**Request:**

```http
POST /2fa/enableOtp
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Success Response (200):**

```json
{
  "message": "2FA enabled successfully"
}
```

---

### 10. Send OTP

**Request:**

```http
POST /2fa/sendOtp
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "email": "john@example.com"
}
```

**Success Response (200):**

```json
{
  "message": "OTP sent to your email"
}
```

---

### 11. Verify 2FA

**Request:**

```http
POST /auth/2fa/verify
Content-Type: application/json

{
  "email": "john@example.com",
  "code": "123456"
}
```

**Success Response (200):**

```json
{
  "message": "2FA verification successful"
}
```

---

### 12. Disable 2FA

**Request:**

```http
POST /2fa/disableOtp
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Success Response (200):**

```json
{
  "message": "2FA disabled successfully"
}
```

---

## User Management

### 13. Get Profile

**Request:**

```http
GET /user/getProfile
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Success Response (200):**

```json
{
  "id": 1,
  "first_name": "John",
  "last_name": "Doe",
  "email": "john@example.com",
  "is_active": true,
  "created_at": "2025-11-27T10:00:00Z",
  "updated_at": "2025-11-27T10:00:00Z"
}
```

**Error Response (401):**

```json
{
  "error": "Unauthorized"
}
```

---

### 14. Update Profile

**Request:**

```http
PUT /user/updateProfile
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "first_name": "Jonathan",
  "last_name": "Smith",
  "email": "jonathan@example.com"
}
```

**Success Response (200):**

```json
{
  "id": 1,
  "first_name": "Jonathan",
  "last_name": "Smith",
  "email": "jonathan@example.com",
  "is_active": true,
  "created_at": "2025-11-27T10:00:00Z",
  "updated_at": "2025-11-27T11:30:00Z"
}
```

---

## Health Check

### 15. Health Status

**Request:**

```http
GET /health
```

**Success Response (200):**

```json
{
  "status": "ok"
}
```

---

## Error Codes

| Code | Status            | Description                          |
| ---- | ----------------- | ------------------------------------ |
| 400  | Bad Request       | Invalid input, validation failed     |
| 401  | Unauthorized      | Missing or invalid token/credentials |
| 404  | Not Found         | Resource not found                   |
| 409  | Conflict          | Email already exists                 |
| 429  | Too Many Requests | Rate limit exceeded                  |
| 500  | Server Error      | Internal server error                |

### Error Response Format

```json
{
  "error": "Error message",
  "validation_error": {
    "field_name": "Error description"
  }
}
```

---

## Password Requirements

- ✓ Minimum 8 characters
- ✓ At least one uppercase letter (A-Z)
- ✓ At least one lowercase letter (a-z)
- ✓ At least one digit (0-9)
- ✓ At least one special character (!@#$%^&\*-\_=+[]{}|;:',.<>?/\`~)

---

## Environment Configuration

Create `.env` file in the root directory:

```env
# =============== APP CONFIG ==================
APP_NAME=Authentio
APP_VERSION=1.0.0
APP_ENV=production
SERVER_PORT=8080
BASE_URL=https://yourdomain.com
FRONTEND_URL=https://app.yourdomain.com

# =============== SECURITY ====================
JWT_SECRET=generate-strong-random-key-min-32-chars
ACCESS_TOKEN_TTL=15m
REFRESH_TOKEN_TTL=168h
BCRYPT_COST=12
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=1m

# =============== DATABASE ====================
POSTGRES_DSN=postgresql://user:password@host:5432/authentio
REDIS_ADDR=redis-host:6379
REDIS_PASS=redis-password

# =============== GOOGLE OAUTH ================
GOOGLE_CLIENT_ID=your-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-secret
GOOGLE_REDIRECT_URL=https://yourdomain.com/api/v1/auth/google/callback

# =============== EMAIL =======================
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=app-specific-password
SMTP_FROM=noreply@yourdomain.com

# =============== LOGGING =====================
LOG_LEVEL=info
ENABLE_REQUEST_LOGS=true

# =============== CORS =======================
CORS_ALLOWED_ORIGINS=https://app.yourdomain.com,https://yourdomain.com
```

---

## Testing with cURL

### Register User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "test@example.com",
    "password": "SecurePass123!"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "SecurePass123!"
  }'
```

### Get Profile (Replace TOKEN with actual access token)

```bash
curl -X GET http://localhost:8080/api/v1/user/getProfile \
  -H "Authorization: Bearer TOKEN"
```

---

## Deployment

### Using Docker

```bash
docker-compose up -d
```

### Using Kubernetes

Refer to `infra/` directory for Kubernetes manifests

### Production Checklist

- [ ] Set `APP_ENV=production` in `.env`
- [ ] Use strong JWT_SECRET (min 32 random chars)
- [ ] Configure real SMTP server (Gmail, SendGrid, etc.)
- [ ] Set up Google OAuth2 credentials
- [ ] Enable HTTPS/TLS
- [ ] Configure CORS for your frontend domain
- [ ] Set up PostgreSQL backups
- [ ] Configure Redis for production
- [ ] Set up monitoring and logging
- [ ] Implement rate limiting rules

---

## Support

For issues and questions:

- ��� Email: support@authentio.com
- ��� GitHub Issues: [GitHub Repository]
- ��� Documentation: [Full Documentation]

---

## License

MIT License - See LICENSE file for details

---

**Built with ❤️ using Go**

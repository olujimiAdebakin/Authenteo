# Authentio API

## Overview
Authentio is a robust authentication and user management API service built with Go and the Gin framework. It provides a complete backend solution for handling user registration, login, profile management, and two-factor authentication, utilizing PostgreSQL for data persistence and JWT for secure sessions.

## Features
- **Go (Golang)**: Core language for building a high-performance, concurrent backend.
- **Gin**: A lightweight and fast HTTP web framework for routing and middleware.
- **PostgreSQL**: The primary relational database for storing user credentials and profile data.
- **JWT (JSON Web Tokens)**: Used for stateless, secure authentication between the client and server.
- **Docker & Docker Compose**: For containerizing the application and its dependencies (Postgres, Redis) for reproducible development and deployment environments.
- **Two-Factor Authentication (2FA)**: OTP-based 2FA flow for enhanced security.

## Getting Started
### Installation
This project uses Docker Compose to streamline the setup process for the application, database, and other services.

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/olujimiAdebakin/Authente-.git
    cd Authente-
    ```

2.  **Create an environment file:**
    Create a `.env` file in the root directory and populate it with the necessary environment variables listed below.

3.  **Build and run the containers:**
    ```bash
    docker-compose up --build
    ```
    The API server will be available at `http://localhost:8080`.

### Environment Variables
Create a `.env` file in the project root and add the following variables.

```env
# Server Configuration
SERVER_PORT=8080
APP_ENV=development

# Database Configuration (matches docker-compose.yml)
POSTGRES_DSN=postgres://postgres:secret@db:5432/authentio_db?sslmode=disable

# Redis Configuration (matches docker-compose.yml)
REDIS_ADDR=redis:6379
REDIS_PASS=

# JWT Configuration
JWT_SECRET=a_very_strong_and_long_secret_key
ACCESS_TOKEN_TTL=15m
REFRESH_TOKEN_TTL=168h

# SMTP Configuration for sending emails (e.g., OTPs)
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=user@example.com
SMTP_PASSWORD=your_smtp_password
```

## API Documentation
### Base URL
All API endpoints are prefixed with `/api/v1`.

### Endpoints
#### POST /auth/register
Registers a new user in the system.

**Request**:
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "password": "a-strong-password123"
}
```

**Response**: (201 Created)
```json
{
    "user": {
        "id": 1,
        "first_name": "John",
        "last_name": "Doe",
        "email": "john.doe@example.com",
        "is_active": true
    },
    "message": "Registration successful"
}
```

**Errors**:
- `400 Bad Request`: Validation failed (e.g., invalid email, short password) or the email already exists.

#### POST /auth/login
Authenticates a user and returns a JWT access token.

**Request**:
```json
{
  "email": "john.doe@example.com",
  "password": "a-strong-password123"
}
```

**Response**: (200 OK)
```json
{
    "user": {
        "id": 1,
        "first_name": "John",
        "last_name": "Doe",
        "email": "john.doe@example.com",
        "is_active": true
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Errors**:
- `401 Unauthorized`: Invalid email or password.

#### POST /auth/refresh
Refreshes an expired access token using a valid refresh token.
*(Note: Handler implementation for this route is not present in the provided source code)*

**Request**:
```json
{
  "refresh_token": "your_valid_refresh_token_here"
}
```

**Response**: (200 OK)
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Errors**:
- `401 Unauthorized`: The refresh token is invalid, expired, or revoked.

#### POST /auth/2fa/send
Sends a One-Time Password (OTP) to the user's registered email for 2FA verification.

**Request**:
```json
{
  "email": "john.doe@example.com"
}
```

**Response**: (200 OK)
```json
{
  "message": "OTP sent successfully"
}
```

**Errors**:
- `400 Bad Request`: Invalid email format.
- `500 Internal Server Error`: Failed to send OTP, or user not found.

#### POST /auth/2fa/verify
Verifies a 2FA OTP code provided by the user.

**Request**:
```json
{
  "email": "john.doe@example.com",
  "code": "123456"
}
```

**Response**: (200 OK)
```json
{
  "message": "OTP verified successfully"
}
```

**Errors**:
- `400 Bad Request`: The provided code is invalid, expired, or does not match.

#### GET /user/me
Retrieves the profile of the currently authenticated user. Requires `Authorization: Bearer <token>` header.

**Request**:
(No request body)

**Response**: (200 OK)
```json
{
    "id": 1,
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com",
    "is_active": true
}
```

**Errors**:
- `401 Unauthorized`: Token is missing, invalid, or expired.
- `500 Internal Server Error`: User associated with the token could not be found.

#### PUT /user/update
Updates the profile of the currently authenticated user. Requires `Authorization: Bearer <token>` header.

**Request**:
```json
{
  "first_name": "Jonathan",
  "last_name": "Doe",
  "email": "jonathan.doe@example.com"
}
```

**Response**: (200 OK)
```json
{
  "message": "Profile updated successfully"
}
```

**Errors**:
- `400 Bad Request`: Invalid input data (e.g., invalid email format).
- `401 Unauthorized`: Token is missing, invalid, or expired.
- `500 Internal Server Error`: An error occurred during the update process (e.g., email already taken).

#### GET /health
A health check endpoint to verify that the service is running.

**Request**:
(No request body)

**Response**: (200 OK)
```json
{
  "status": "ok"
}
```

**Errors**:
- None.
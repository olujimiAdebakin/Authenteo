# Authentio API

## Overview
Authentio is a secure, high-performance authentication and authorization API built with Go (Golang) and the Gin web framework. It provides a complete solution for user management, including JWT-based sessions, Google OAuth2, email-based Two-Factor Authentication (2FA), and secure password reset flows, all backed by PostgreSQL and Redis.

## Features
- **Go (Golang) & Gin**: A robust and efficient backend for handling high-concurrency API requests.
- **PostgreSQL**: The primary relational database for persistent user and token storage.
- **Redis**: Used for distributed rate limiting and token blacklisting to enhance security and performance.
- **JWT Authentication**: Stateless and secure authentication using JSON Web Tokens.
- **Google OAuth2**: Seamless social login integration for a frictionless user experience.
- **Two-Factor Authentication (2FA)**: Email-based OTP verification for an added layer of security.
- **Docker & Docker Compose**: Fully containerized for easy setup, deployment, and scalability.

## Getting Started
The recommended way to run this project is using Docker.

### Installation
1.  **Clone the repository:**
    ```bash
    git clone https://github.com/your-username/authentio.git
    cd authentio
    ```

2.  **Create an environment file:**
    Create a `.env` file in the root directory and populate it with the variables listed below.

3.  **Build and run with Docker Compose:**
    ```bash
    docker-compose up --build -d
    ```
    The API will be available at `http://localhost:8080`.

### Environment Variables
Create a `.env` file in the project root with the following variables:

| Variable                | Description                                                | Example                                                        |
| ----------------------- | ---------------------------------------------------------- | -------------------------------------------------------------- |
| `SERVER_PORT`           | Port for the API server.                                   | `8080`                                                         |
| `APP_ENV`               | Application environment.                                   | `development` or `production`                                  |
| `POSTGRES_DSN`          | PostgreSQL connection string.                              | `postgres://postgres:secret@authentio_db:5432/authentio_db?sslmode=disable` |
| `REDIS_ADDR`            | Redis server address.                                      | `authentio_redis:6379`                                         |
| `REDIS_PASS`            | Redis password (if any).                                   | ``                                                             |
| `JWT_SECRET`            | Secret key for signing JWTs.                               | `a_very_strong_and_long_secret_key`                            |
| `SMTP_HOST`             | SMTP server for sending emails.                            | `smtp.gmail.com`                                               |
| `SMTP_PORT`             | SMTP server port.                                          | `587`                                                          |
| `SMTP_USERNAME`         | Your email address for sending emails.                     | `your.email@example.com`                                       |
| `SMTP_PASSWORD`         | Your email app password or token.                          | `yourapppassword`                                              |
| `SMTP_FROM`             | The "From" address for outgoing emails.                    | `your.email@example.com`                                       |
| `GOOGLE_CLIENT_ID`      | Google OAuth2 Client ID.                                   | `your-google-client-id.apps.googleusercontent.com`             |
| `GOOGLE_CLIENT_SECRET`  | Google OAuth2 Client Secret.                               | `GOCSPX-your-secret`                                           |
| `GOOGLE_REDIRECT_URL`   | Callback URL registered with Google.                       | `http://localhost:8080/api/v1/auth/google/callback`            |

## API Documentation
### Base URL
`/api/v1`

### Endpoints
#### POST /auth/register
**Description**: Registers a new user.

**Request**:
```json
{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com",
    "password": "StrongPassword123!"
}
```

**Response** (201 Created):
```json
{
    "user": {
        "id": 1,
        "first_name": "John",
        "last_name": "Doe",
        "email": "john.doe@example.com",
        "is_active": true,
        "created_at": "2023-10-27T10:00:00Z"
    },
    "message": "Registration successful"
}
```

**Errors**:
- `400 Bad Request`: Validation failed (e.g., weak password, invalid email).
- `409 Conflict`: Email already exists.

---
#### POST /auth/login
**Description**: Authenticates a user and returns JWT tokens.

**Request**:
```json
{
    "email": "john.doe@example.com",
    "password": "StrongPassword123!"
}
```

**Response** (200 OK):
```json
{
    "user": {
        "id": 1,
        "first_name": "John",
        "last_name": "Doe",
        "email": "john.doe@example.com",
        "is_active": true
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "a1b2c3d4e5f6...",
    "expires_in": 3600
}
```

**Errors**:
- `401 Unauthorized`: Invalid email or password.

---
#### POST /auth/refresh
**Description**: Generates a new access token using a valid refresh token.

**Request**:
```json
{
    "refresh_token": "a1b2c3d4e5f6..."
}
```

**Response** (200 OK):
```json
{
    "user": {
        "id": 1,
        "first_name": "John",
        "last_name": "Doe",
        "email": "john.doe@example.com",
        "is_active": true
    },
    "access_token": "new_eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "new_a1b2c3d4e5f6...",
    "expires_in": 3600
}
```

**Errors**:
- `400 Bad Request`: Invalid or expired refresh token.

---
#### POST /auth/google/login
**Description**: Authenticates a user via a Google ID token (client-side flow).

**Request**:
```json
{
    "id_token": "google_id_token_from_frontend"
}
```

**Response** (200 OK):
```json
{
    "user": {
        "id": 2,
        "first_name": "Jane",
        "last_name": "Doe",
        "email": "jane.doe.google@example.com",
        "is_active": true
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "a1b2c3d4e5f6...",
    "expires_in": 3600
}
```

**Errors**:
- `401 Unauthorized`: Invalid Google token.

---
#### POST /auth/forgot-password
**Description**: Sends a password reset code to the user's email.

**Request**:
```json
{
    "email": "john.doe@example.com"
}
```

**Response** (200 OK):
```json
{
    "message": "Password reset email sent"
}
```

**Errors**:
- `400 Bad Request`: Invalid email format.

---
#### POST /auth/reset-password
**Description**: Resets the user's password using the code from the email.

**Request**:
```json
{
    "email": "john.doe@example.com",
    "code": "123456",
    "new_password": "NewStrongPassword123!"
}
```

**Response** (200 OK):
```json
{
    "message": "Password reset successful"
}
```

**Errors**:
- `400 Bad Request`: Invalid code, expired code, or weak new password.

---
#### POST /auth/2fa/verify
**Description**: Verifies a 2FA code during the login process for users with 2FA enabled.

**Request**:
```json
{
    "email": "john.doe@example.com",
    "code": "654321"
}
```

**Response** (200 OK):
```json
{
    "message": "2FA verification successful"
}
```

**Errors**:
- `400 Bad Request`: Invalid or expired 2FA code.

---
#### POST /2fa/enableOtp
**Description**: Enables email-based 2FA for the authenticated user. (Requires Bearer Token)

**Request**:
*(No request body)*

**Response** (200 OK):
```json
{
    "message": "2FA enabled successfully"
}
```

**Errors**:
- `401 Unauthorized`: Invalid or missing JWT token.

---
#### POST /2fa/disableOtp
**Description**: Disables 2FA for the authenticated user. (Requires Bearer Token)

**Request**:
*(No request body)*

**Response** (200 OK):
```json
{
    "message": "2FA disabled successfully"
}
```

**Errors**:
- `401 Unauthorized`: Invalid or missing JWT token.

---
#### POST /2fa/sendOtp
**Description**: Sends a new 2FA OTP code to the user's email. (Requires Bearer Token)

**Request**:
```json
{
    "email": "john.doe@example.com"
}
```

**Response** (200 OK):
```json
{
    "message": "OTP sent successfully"
}
```

**Errors**:
- `401 Unauthorized`: Invalid or missing JWT token.
- `500 Internal Server Error`: Failed to send OTP email.

---
#### GET /user/getProfile
**Description**: Retrieves the authenticated user's profile. (Requires Bearer Token)

**Request**:
*(No request body)*

**Response** (200 OK):
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
- `401 Unauthorized`: Invalid or missing JWT token.
- `404 Not Found`: User not found.

---
#### PUT /user/updateProfile
**Description**: Updates the authenticated user's profile information. (Requires Bearer Token)

**Request**:
```json
{
    "first_name": "Jonathan",
    "last_name": "Doer",
    "email": "jonathan.doer@example.com"
}
```

**Response** (200 OK):
```json
{
    "message": "Profile updated successfully"
}
```

**Errors**:
- `401 Unauthorized`: Invalid or missing JWT token.
- `409 Conflict`: New email is already in use.

---
#### GET /health
**Description**: A public endpoint to check the health of the service.

**Request**:
*(No request body)*

**Response** (200 OK):
```json
{
    "status": "ok"
}
```

**Errors**:
- None.

[![Readme was generated by Dokugen](https://img.shields.io/badge/Readme%20was%20generated%20by-Dokugen-brightgreen)](https://www.npmjs.com/package/dokugen)
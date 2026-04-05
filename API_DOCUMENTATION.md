# API Documentation — Learn Fiber Auth JWT

Base URL: `http://localhost:{PORT}/api/v1`

## Table of Contents

- [Overview](#overview)
- [Authentication](#authentication)
- [Common Response Format](#common-response-format)
- [Environment Variables](#environment-variables)
- [Database Schema](#database-schema)
- [Endpoints](#endpoints)
  - [Health Check](#health-check)
  - [Register](#register)
  - [Login](#login)
  - [Get User Details](#get-user-details)
  - [Revoke Access (Logout)](#revoke-access-logout)
  - [Verify Email](#verify-email)
- [Error Codes](#error-codes)

---

## Overview

A RESTful authentication API built with [Go Fiber](https://gofiber.io/) featuring:

- User registration and login with JWT-based authentication
- Access token (JWT, 1-hour expiry) + Refresh token (UUID, 7-day expiry)
- Email verification via SMTP
- Password hashing with bcrypt
- PostgreSQL database with UUID primary keys

### CORS Configuration

| Setting            | Value                                            |
| ------------------ | ------------------------------------------------ |
| Allowed Origins    | `http://localhost:3000`, `http://localhost:5173`  |
| Allowed Methods    | GET, POST, PUT, DELETE, OPTIONS                  |
| Allowed Headers    | Origin, Content-Type, Authorization              |
| Exposed Headers    | Content-Length                                    |
| Allow Credentials  | true                                             |
| Max Age            | 12 hours                                         |

---

## Authentication

Protected endpoints require a valid JWT access token in the `Authorization` header:

```
Authorization: Bearer <access_token>
```

### Token Details

| Token Type    | Format      | Signing Algorithm | Expiration | Storage        |
| ------------- | ----------- | ----------------- | ---------- | -------------- |
| Access Token  | JWT         | HS256             | 1 hour     | Client-side    |
| Refresh Token | UUID string | —                 | 7 days     | Database       |

**JWT Claims:**

| Claim     | Description                       |
| --------- | --------------------------------- |
| `auth_id` | User UUID                         |
| `iss`     | `learn-fiber-auth-jwt`            |
| `sub`     | User UUID (string)                |
| `iat`     | Issued at (Unix timestamp)        |
| `nbf`     | Not before (Unix timestamp)       |
| `exp`     | Expiration time (Unix timestamp)  |

---

## Common Response Format

All responses follow a generic wrapper structure:

### Success Response

```json
{
  "data": { ... },
  "status": "SUCCESS",
  "error": null
}
```

### Error Response

```json
{
  "data": null,
  "status": "FAIL",
  "error": {
    "http_code": 400,
    "error_code": "ERROR-CODE",
    "message": "Error description"
  }
}
```

### Status Values

| Status    | Description                          |
| --------- | ------------------------------------ |
| `SUCCESS` | Request completed successfully       |
| `FAIL`    | Request failed due to client error   |
| `ERROR`   | Request failed due to server error   |

---

## Environment Variables

| Variable      | Description                    | Example              |
| ------------- | ------------------------------ | -------------------- |
| `PORT`        | Server port                    | `8003`               |
| `DB_HOST`     | PostgreSQL host                | `localhost`          |
| `DB_PORT`     | PostgreSQL port                | `5432`               |
| `DB_USER`     | PostgreSQL user                | `postgres`           |
| `DB_PASSWORD` | PostgreSQL password            | `password`           |
| `DB_NAME`     | PostgreSQL database name       | `auth_db`            |
| `SECRET`      | JWT signing secret key (HS256) | `my-secret-key`      |
| `SMTP_HOST`   | SMTP server host               | `smtp.gmail.com`     |
| `SMTP_PORT`   | SMTP server port               | `587`                |
| `SMTP_USER`   | SMTP username                  | `user@gmail.com`     |
| `SMTP_PASS`   | SMTP password                  | `app-password`       |

---

## Database Schema

### `tb_auth`

| Column        | Type           | Constraints                       |
| ------------- | -------------- | --------------------------------- |
| `id`          | UUID           | PRIMARY KEY, DEFAULT gen_random_uuid() |
| `created_at`  | TIMESTAMP      | NOT NULL, DEFAULT NOW()           |
| `updated_at`  | TIMESTAMP      | NOT NULL, DEFAULT NOW()           |
| `deleted_at`  | TIMESTAMP      | Nullable (soft delete)            |
| `username`    | VARCHAR(255)   | NOT NULL, UNIQUE                  |
| `email`       | VARCHAR(255)   | NOT NULL, UNIQUE                  |
| `password`    | VARCHAR(255)   | NOT NULL (bcrypt hash)            |
| `is_verified` | BOOLEAN        | NOT NULL, DEFAULT false           |

### `tb_refresh_tokens`

| Column        | Type           | Constraints                       |
| ------------- | -------------- | --------------------------------- |
| `id`          | UUID           | PRIMARY KEY, DEFAULT gen_random_uuid() |
| `created_at`  | TIMESTAMP      | NOT NULL, DEFAULT NOW()           |
| `updated_at`  | TIMESTAMP      | NOT NULL, DEFAULT NOW()           |
| `deleted_at`  | TIMESTAMP      | Nullable (soft delete)            |
| `auth_id`     | UUID           | NOT NULL, FK → tb_auth(id) ON DELETE CASCADE |
| `token`       | VARCHAR(255)   | NOT NULL, UNIQUE                  |
| `expires_at`  | TIMESTAMP      | NOT NULL                          |
| `revoked`     | BOOLEAN        | NOT NULL, DEFAULT false           |

### `tb_email_verify_tokens`

| Column        | Type           | Constraints                       |
| ------------- | -------------- | --------------------------------- |
| `id`          | UUID           | PRIMARY KEY, DEFAULT gen_random_uuid() |
| `created_at`  | TIMESTAMP      | NOT NULL, DEFAULT NOW()           |
| `updated_at`  | TIMESTAMP      | NOT NULL, DEFAULT NOW()           |
| `deleted_at`  | TIMESTAMP      | Nullable (soft delete)            |
| `auth_id`     | UUID           | NOT NULL, UNIQUE, FK → tb_auth(id) ON DELETE CASCADE |
| `token`       | VARCHAR(255)   | NOT NULL, UNIQUE                  |
| `expires_at`  | TIMESTAMP      | NOT NULL                          |

---

## Endpoints

---

### Health Check

Check if the API server is running.

**`GET /api/v1/health-check/health`**

**Authentication:** None

#### Response

**`200 OK`**

```json
{
  "data": "OK",
  "status": "SUCCESS",
  "error": null
}
```

---

### Register

Register a new user account. A verification email is sent asynchronously after successful registration.

**`POST /api/v1/auth/register`**

**Authentication:** None

#### Request Body

| Field      | Type   | Required | Validation        | Description         |
| ---------- | ------ | -------- | ----------------- | ------------------- |
| `username` | string | Yes      | `required`        | Unique username     |
| `password` | string | Yes      | `required`        | Plain text password |
| `email`    | string | Yes      | `required, email` | Valid email address |

```json
{
  "username": "johndoe",
  "password": "securepassword123",
  "email": "johndoe@example.com"
}
```

#### Responses

**`201 Created`** — Registration successful

```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "johndoe",
    "email": "johndoe@example.com",
    "is_verified": false
  },
  "status": "SUCCESS",
  "error": null
}
```

**`400 Bad Request`** — Validation error (missing or invalid fields)

```json
{
  "data": null,
  "status": "FAIL",
  "error": {
    "http_code": 400,
    "error_code": "VALIDATE-01",
    "message": "Validation error details per field"
  }
}
```

**`400 Bad Request`** — Registration failed (e.g., duplicate username or email)

```json
{
  "data": null,
  "status": "FAIL",
  "error": {
    "http_code": 400,
    "error_code": "REGISTER-01",
    "message": "Error description"
  }
}
```

#### Side Effects

- Password is hashed with bcrypt (default cost 10)
- A 64-character hex email verification token is generated (valid for 24 hours)
- A verification email is sent to the registered email address with a link:
  `http://localhost:8003/api/v1/auth/verify-email?token={token}`
- Email sending failures are logged but do not fail the registration

---

### Login

Authenticate a user and receive access + refresh tokens.

**`POST /api/v1/auth/login`**

**Authentication:** None

#### Request Body

| Field      | Type   | Required | Validation | Description         |
| ---------- | ------ | -------- | ---------- | ------------------- |
| `username` | string | Yes      | `required` | Registered username |
| `password` | string | Yes      | `required` | User password       |

```json
{
  "username": "johndoe",
  "password": "securepassword123"
}
```

#### Responses

**`200 OK`** — Login successful

```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "550e8400-e29b-41d4-a716-446655440000"
  },
  "status": "SUCCESS",
  "error": null
}
```

**`400 Bad Request`** — Validation error

```json
{
  "data": null,
  "status": "FAIL",
  "error": {
    "http_code": 400,
    "error_code": "VALIDATE-01",
    "message": "Validation error details"
  }
}
```

**`400 Bad Request`** — Invalid credentials

```json
{
  "data": null,
  "status": "FAIL",
  "error": {
    "http_code": 400,
    "error_code": "LOGIN-01",
    "message": "Invalid username or password"
  }
}
```

#### Side Effects

- A new refresh token (UUID) is stored in the database with 7-day expiry
- Access token (JWT) is generated with 1-hour expiry

---

### Get User Details

Retrieve the authenticated user's profile information.

**`GET /api/v1/auth/detail`**

**Authentication:** Required (Bearer Token)

#### Headers

| Header          | Value                  | Required |
| --------------- | ---------------------- | -------- |
| `Authorization` | `Bearer <access_token>` | Yes      |

#### Responses

**`200 OK`** — User details retrieved

```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "johndoe",
    "email": "johndoe@example.com",
    "is_verified": true
  },
  "status": "SUCCESS",
  "error": null
}
```

**`401 Unauthorized`** — Missing, invalid, or expired token

```json
{
  "data": null,
  "status": "FAIL",
  "error": {
    "http_code": 401,
    "error_code": "AUTH-01",
    "message": "Unauthorized"
  }
}
```

**`500 Internal Server Error`** — Server error

```json
{
  "data": null,
  "status": "ERROR",
  "error": {
    "http_code": 500,
    "error_code": "DETAIL-01",
    "message": "Error description"
  }
}
```

---

### Revoke Access (Logout)

Revoke a refresh token to effectively log out a user.

**`POST /api/v1/auth/revoke`**

**Authentication:** None

#### Request Body

| Field           | Type   | Required | Description                   |
| --------------- | ------ | -------- | ----------------------------- |
| `refresh_token` | string | Yes      | The refresh token to revoke   |

```json
{
  "refresh_token": "550e8400-e29b-41d4-a716-446655440000"
}
```

#### Responses

**`200 OK`** — Token revoked successfully

```json
{
  "data": "Access revoked",
  "status": "SUCCESS",
  "error": null
}
```

**`400 Bad Request`** — Invalid request body

```json
{
  "data": null,
  "status": "FAIL",
  "error": {
    "http_code": 400,
    "error_code": "REVOKE-01",
    "message": "Invalid request body"
  }
}
```

**`500 Internal Server Error`** — Server error during revocation

```json
{
  "data": null,
  "status": "ERROR",
  "error": {
    "http_code": 500,
    "error_code": "REVOKE-02",
    "message": "Error description"
  }
}
```

#### Side Effects

- The refresh token's `revoked` flag is set to `true` in the database
- The `updated_at` timestamp is updated

---

### Verify Email

Verify a user's email address using the token received via email.

**`GET /api/v1/auth/verify-email`**

**Authentication:** None

#### Query Parameters

| Parameter | Type   | Required | Description                                          |
| --------- | ------ | -------- | ---------------------------------------------------- |
| `token`   | string | Yes      | 64-character hex verification token from email link   |

#### Example Request

```
GET /api/v1/auth/verify-email?token=a1b2c3d4e5f6...
```

#### Responses

**`200 OK`** — Email verified successfully

```json
{
  "data": "Email verified successfully",
  "status": "SUCCESS",
  "error": null
}
```

**`400 Bad Request`** — Missing token parameter

```json
{
  "data": null,
  "status": "FAIL",
  "error": {
    "http_code": 400,
    "error_code": "VERIFY-01",
    "message": "Token is required"
  }
}
```

**`400 Bad Request`** — Invalid or expired token

```json
{
  "data": null,
  "status": "FAIL",
  "error": {
    "http_code": 400,
    "error_code": "VERIFY-02",
    "message": "Invalid or expired token"
  }
}
```

#### Side Effects

- User's `is_verified` field is set to `true`
- The verification token is deleted from the database after use

---

## Error Codes

| Error Code    | HTTP Status | Endpoint       | Description                                  |
| ------------- | ----------- | -------------- | -------------------------------------------- |
| `VALIDATE-01` | 400         | Multiple       | Request body validation failed               |
| `REGISTER-01` | 400         | Register       | Registration failed (duplicate user/email)   |
| `LOGIN-01`    | 400         | Login          | Invalid username or password                 |
| `AUTH-01`     | 401         | Protected      | Missing, invalid, or expired JWT token       |
| `DETAIL-01`   | 500         | Get Details    | Failed to retrieve user details              |
| `REVOKE-01`   | 400         | Revoke Access  | Invalid request body                         |
| `REVOKE-02`   | 500         | Revoke Access  | Server error during token revocation         |
| `VERIFY-01`   | 400         | Verify Email   | Missing token query parameter                |
| `VERIFY-02`   | 400         | Verify Email   | Invalid or expired verification token        |

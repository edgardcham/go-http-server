# Go HTTP Server - Chirpy

A learning project for building a RESTful API server in Go. This project implements a Twitter-like service called "Chirpy" where users can create accounts, post short messages (chirps), and interact with the platform.

## Features

- User authentication with JWT tokens
- Password hashing with bcrypt
- Refresh token system for extended sessions
- CRUD operations for chirps (tweets)
- Premium user upgrades via webhooks
- Database migrations with Goose
- SQL query generation with sqlc

## Tech Stack

- **Language**: Go 1.24+
- **Database**: PostgreSQL
- **Authentication**: JWT (golang-jwt)
- **Password Hashing**: bcrypt
- **Database Migrations**: Goose
- **SQL Generation**: sqlc
- **Environment Variables**: godotenv

## Project Structure

```
.
├── main.go                 # Entry point and server setup
├── responses.go            # Response helper functions
├── handler_*.go           # HTTP handlers for each endpoint
├── middleware_metrics.go   # Metrics middleware
├── internal/
│   ├── auth/              # Authentication utilities
│   │   ├── auth.go        # Password hashing functions
│   │   └── jwt.go         # JWT creation/validation
│   └── database/          # Generated database code (sqlc)
├── sql/
│   ├── schema/            # Database migrations
│   └── queries/           # SQL queries for sqlc
└── vendor/                # Dependencies
```

## API Endpoints

### Public Endpoints

#### Health Check

```
GET /api/healthz
```

Returns server status.

**Response**: `200 OK`

```json
{
  "status": "ok"
}
```

### User Management

#### Create User

```
POST /api/users
```

Creates a new user account.

**Request Body**:

```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response**: `201 Created`

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z",
  "is_chirpy_red": false
}
```

#### User Login

```
POST /api/login
```

Authenticates user and returns access + refresh tokens.

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "securepassword123",
  "expires_in_seconds": 3600  // optional, defaults to 1 hour
}
```

**Response**: `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z",
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "e3713674b1112350566c42137bf1ae485e6751356e741efe20c6d209c23df6c7",
  "is_chirpy_red": false
}
```

#### Update User
```
PUT /api/users
Authorization: Bearer <access_token>
```
Updates user email and/or password.

**Request Body**:
```json
{
  "email": "newemail@example.com",
  "password": "newpassword123"
}
```

**Response**: `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "newemail@example.com",
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:30:00Z",
  "is_chirpy_red": false
}
```

### Token Management

#### Refresh Access Token
```
POST /api/refresh
Authorization: Bearer <refresh_token>
```
Exchanges a valid refresh token for a new access token.

**Response**: `200 OK`
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

#### Revoke Refresh Token
```
POST /api/revoke
Authorization: Bearer <refresh_token>
```
Revokes a refresh token (logout).

**Response**: `204 No Content`

### Chirps (Posts)

#### Create Chirp
```
POST /api/chirps
Authorization: Bearer <access_token>
```
Creates a new chirp for authenticated user.

**Request Body**:
```json
{
  "body": "This is my chirp!"
}
```

**Response**: `201 Created`
```json
{
  "id": "650e8400-e29b-41d4-a716-446655440000",
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z",
  "body": "This is my chirp!",
  "user_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Note**: Chirps are limited to 140 characters. Profane words (kerfuffle, sharbert, fornax) are automatically replaced with ****.

#### Get All Chirps
```
GET /api/chirps?author_id=<user_id>&sort=<asc|desc>
```
Retrieves all chirps, optionally filtered by author.

**Query Parameters**:
- `author_id` (optional): Filter chirps by user ID
- `sort` (optional): Sort by created_at. Values: `asc` (default) or `desc`

**Response**: `200 OK`
```json
[
  {
    "id": "650e8400-e29b-41d4-a716-446655440000",
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z",
    "body": "This is my chirp!",
    "user_id": "550e8400-e29b-41d4-a716-446655440000"
  }
]
```

#### Get Single Chirp
```
GET /api/chirps/{chirpID}
```
Retrieves a specific chirp by ID.

**Response**: `200 OK`
```json
{
  "id": "650e8400-e29b-41d4-a716-446655440000",
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z",
  "body": "This is my chirp!",
  "user_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

#### Delete Chirp
```
DELETE /api/chirps/{chirpID}
Authorization: Bearer <access_token>
```
Deletes a chirp (only by the chirp's author).

**Response**: `204 No Content`

### Premium Features

#### Polka Webhook
```
POST /api/polka/webhooks
Authorization: ApiKey <polka_api_key>
```
Webhook endpoint for Polka payment system to upgrade users to Chirpy Red.

**Request Body**:
```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

**Response**: `204 No Content`

### Admin Endpoints

#### Metrics
```
GET /admin/metrics
```
Displays server metrics (development only).

**Response**: `200 OK` (HTML)

#### Reset Database
```
POST /admin/reset
```
Resets the database and metrics (development only).

**Response**: `200 OK`
```
Hits reset to 0
```

## Authentication

The API uses JWT tokens for authentication:

1. **Access Tokens**: Short-lived (1 hour default), used for API requests
2. **Refresh Tokens**: Long-lived (60 days), used to obtain new access tokens

Include the access token in the Authorization header:
```
Authorization: Bearer <access_token>
```

## Environment Variables

Create a `.env` file in the project root:

```env
DB_URL=postgres://username:password@localhost/chirpy?sslmode=disable
JWT_SECRET_KEY=your-32-byte-secret-key-here
PLATFORM=dev  # or "prod"
POLKA_KEY=your-polka-api-key
```

## Database Setup

1. Install PostgreSQL
2. Create database: `createdb chirpy`
3. Run migrations:
   ```bash
   goose postgres "$DB_URL" up
   ```

## Development

### Running the Server
```bash
go run .
```
Server runs on `http://localhost:8080`

### Database Migrations

Create new migration:
```bash
goose -dir sql/schema create migration_name sql
```

Run migrations:
```bash
goose postgres "$DB_URL" up
```

Rollback:
```bash
goose postgres "$DB_URL" down
```

### Generating Database Code

After modifying SQL queries in `sql/queries/`:
```bash
sqlc generate
```

### Testing

Run the test suite:
```bash
go test ./...
```

## Error Responses

All error responses follow this format:
```json
{
  "error": "Error message here"
}
```

Common status codes:
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Missing or invalid authentication
- `403 Forbidden`: Authenticated but not authorized
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource already exists
- `500 Internal Server Error`: Server error

## Learning Notes

This project demonstrates:
- HTTP server setup with Go's standard library
- Middleware implementation (metrics tracking)
- JWT authentication and refresh token pattern
- Password security with bcrypt
- Database operations with PostgreSQL
- SQL migrations with Goose
- Type-safe SQL with sqlc
- RESTful API design
- Error handling patterns
- Environment configuration
- Webhook integration

## Security Considerations

- Passwords are hashed with bcrypt before storage
- JWT secrets should be long and random (32+ bytes)
- Refresh tokens are stored in database and can be revoked
- HTTPS should be used in production
- API keys for webhooks should be kept secret
- Database connections should use SSL in production
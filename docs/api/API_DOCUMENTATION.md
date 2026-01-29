# Home Control API Documentation

## Overview

Home Control API provides a RESTful interface for managing home automation systems. The API uses JSON for request/response payloads and supports both Bearer token authentication and API key authentication.

## Base URL

```
http://localhost:8080
```

## Authentication

The API supports two authentication methods:

### 1. Bearer Token Authentication

Obtain a token by logging in with username and password, then use it in the `Authorization` header.

**Login Request:**
```bash
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}
```

**Login Response:**
```json
{
  "token": "a1b2c3d4e5f6...",
  "token_type": "bearer",
  "expires_in": 86400,
  "username": "admin",
  "message": "Login successful"
}
```

**Using the Token:**
```bash
GET /api/v1/example
Authorization: Bearer a1b2c3d4e5f6...
```

### 2. API Key Authentication

Use an API key in the `X-API-Key` header. API keys are stored in the SQLite database.

**Request with API Key:**
```bash
GET /api/v1/example
X-API-Key: your-api-key-here
```

## API Endpoints

### Public Endpoints (No Authentication Required)

#### Health Check

**GET** `/health`

Check if the service is running.

**Response:**
```json
{
  "status": "ok",
  "message": "Service is running"
}
```

#### Get API Version

**GET** `/api/v1/version`

Get information about the API version.

**Response:**
```json
{
  "version": "0.1.0",
  "name": "home-ctrl"
}
```

### Authentication Endpoints

#### Login

**POST** `/api/v1/auth/login`

Authenticate with username and password to obtain a Bearer token.

**Request Body:**
```json
{
  "username": "string",
  "password": "string"
}
```

**Response:**
```json
{
  "token": "string",
  "token_type": "bearer",
  "expires_in": 86400,
  "username": "string",
  "message": "Login successful"
}
```

**Status Codes:**
- `200 OK` - Login successful
- `400 Bad Request` - Invalid request body
- `401 Unauthorized` - Invalid username or password

#### Logout

**POST** `/api/v1/auth/logout`

Invalidate the current session.

**Headers:**
```
Authorization: Bearer your-token-here
```

**Response:**
```json
{
  "message": "Logout successful"
}
```

**Status Codes:**
- `200 OK` - Logout successful
- `401 Unauthorized` - Invalid or missing token

### Protected Endpoints (Require Authentication)

#### Get Current User Info

**GET** `/api/v1/me`

Get information about the currently authenticated user.

**Headers:**
```
Authorization: Bearer your-token-here
```

**Response:**
```json
{
  "username": "string",
  "message": "You are authenticated!"
}
```

**Status Codes:**
- `200 OK` - Request successful
- `401 Unauthorized` - Invalid or missing authentication

#### Example Protected Endpoint

**GET** `/api/v1/example`

Example endpoint that demonstrates protected access.

**Headers:**
```
Authorization: Bearer your-token-here
```

**Response:**
```json
{
  "message": "Hello from protected endpoint!",
  "user": "username",
  "config": {
    "host": "127.0.0.1",
    "port": 8080
  }
}
```

**Status Codes:**
- `200 OK` - Request successful
- `401 Unauthorized` - Invalid or missing authentication

## Error Responses

The API returns standard HTTP status codes and JSON error responses:

### Common Error Responses

**400 Bad Request**
```json
{
  "error": "Bad Request",
  "message": "Description of the error"
}
```

**401 Unauthorized**
```json
{
  "error": "Unauthorized",
  "message": "Authentication required"
}
```

**404 Not Found**
```json
{
  "error": "Not found",
  "path": "/nonexistent"
}
```

**500 Internal Server Error**
```json
{
  "error": "Internal Server Error",
  "message": "Description of the error"
}
```

## Configuration

The API can be configured using a YAML configuration file:

```yaml
# Server configuration
server:
  host: "0.0.0.0"
  port: 8080

# Authentication configuration
auth:
  users:
    admin: "admin123"
    user: "user123"
  session_ttl_hours: 24

# Data directory for SQLite database
data_dir: "data"
```

## Examples

### Complete Workflow Example

```bash
# 1. Check service health
curl http://localhost:8080/health

# 2. Get API version
curl http://localhost:8080/api/v1/version

# 3. Login to get a token
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}' | jq -r '.token')

# 4. Access protected endpoint
echo "Token: $TOKEN"
curl http://localhost:8080/api/v1/me \
  -H "Authorization: Bearer $TOKEN"

# 5. Access another protected endpoint
echo "Token: $TOKEN"
curl http://localhost:8080/api/v1/example \
  -H "Authorization: Bearer $TOKEN"

# 6. Logout
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer $TOKEN"
```

### Using API Key

```bash
# Get API key from database (default: default-api-key-12345)
API_KEY="default-api-key-12345"

# Access protected endpoint with API key
curl http://localhost:8080/api/v1/example \
  -H "X-API-Key: $API_KEY"
```

## API Key Management

API keys are stored in the SQLite database. You can manage them programmatically:

### Create a new API key

```go
apiKey, err := db.CreateAPIKey("new-key-123", "My API Key", nil)
```

### List all API keys

```go
keys, err := db.ListAPIKeys()
```

### Delete an API key

```go
err := db.DeleteAPIKey("key-to-delete")
```

## Session Management

Sessions are automatically managed by the system:

- Session tokens expire after the configured TTL (default: 24 hours)
- Expired sessions are automatically cleaned up
- Logout explicitly deletes the session

## Rate Limiting

Currently, the API does not implement rate limiting, but this can be added using Gin middleware.

## CORS

CORS is not explicitly configured, but can be added using Gin's CORS middleware if needed for web applications.

## Security

- All sensitive endpoints require authentication
- Passwords are stored in plain text in the configuration (for development only)
- Use HTTPS in production
- Rotate API keys regularly
- Set appropriate session TTL values

## Versioning

The API uses URL versioning with the `/api/v1/` prefix. This allows for backward compatibility when introducing new API versions.

## Future Development

The API is designed to be extensible. Future versions may include:

- Additional authentication methods (OAuth, JWT)
- More granular permissions
- API key rotation
- Rate limiting
- Request/response logging
- OpenAPI/Swagger documentation
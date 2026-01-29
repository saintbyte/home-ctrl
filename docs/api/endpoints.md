# API Endpoints

## Public Endpoints (No Authentication Required)

### Health Check

**GET** `/health`

Check if the service is running.

**Response:**
```json
{
  "status": "ok",
  "message": "Service is running"
}
```

### Get API Version

**GET** `/api/v1/version`

Get information about the API version.

**Response:**
```json
{
  "version": "0.1.0",
  "name": "home-ctrl"
}
```

## Authentication Endpoints

### Login

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

### Logout

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

## Protected Endpoints (Require Authentication)

### Get Current User Info

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

### Example Protected Endpoint

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
# Authentication

## Overview

The Home Control API supports two authentication methods:

1. **Bearer Token Authentication** - for interactive users
2. **API Key Authentication** - for programmatic access

## Bearer Token Authentication

### Login

**Endpoint:** `POST /api/v1/auth/login`

**Request:**
```bash
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}
```

**Response:**
```json
{
  "token": "a1b2c3d4e5f6...",
  "token_type": "bearer",
  "expires_in": 86400,
  "username": "admin",
  "message": "Login successful"
}
```

### Using the Token

Include the token in the `Authorization` header:

```bash
GET /api/v1/example
Authorization: Bearer a1b2c3d4e5f6...
```

### Logout

**Endpoint:** `POST /api/v1/auth/logout`

**Request:**
```bash
POST /api/v1/auth/logout
Authorization: Bearer your-token-here
```

**Response:**
```json
{
  "message": "Logout successful"
}
```

## API Key Authentication

### Using API Keys

Include the API key in the `X-API-Key` header:

```bash
GET /api/v1/example
X-API-Key: your-api-key-here
```

### Managing API Keys

API keys are stored in the SQLite database and can be managed programmatically:

#### Create API Key

```go
apiKey, err := db.CreateAPIKey("new-key-123", "My API Key", nil)
```

#### List API Keys

```go
keys, err := db.ListAPIKeys()
```

#### Delete API Key

```go
err := db.DeleteAPIKey("key-to-delete")
```

## Session Management

- Session tokens expire after the configured TTL (default: 24 hours)
- Expired sessions are automatically cleaned up
- Logout explicitly deletes the session

## Configuration

Authentication can be configured in the YAML configuration file:

```yaml
auth:
  users:
    admin: "admin123"
    user: "user123"
  session_ttl_hours: 24
```

## Security Best Practices

1. **Use HTTPS**: Always use HTTPS in production
2. **Rotate Keys**: Regularly rotate API keys
3. **Strong Passwords**: Use strong passwords for user accounts
4. **Session TTL**: Set appropriate session expiration times
5. **Monitor Access**: Monitor API access logs for suspicious activity
# API Examples

## Complete Workflow Example

This example demonstrates a complete workflow from login to accessing protected endpoints.

```bash
# 1. Check service health
curl http://localhost:8080/health

# 2. Get API version
curl http://localhost:8080/api/v1/version

# 3. Login to get a token
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}' | jq -r '.token')

echo "Received token: $TOKEN"

# 4. Access protected endpoint - get user info
curl http://localhost:8080/api/v1/me \
  -H "Authorization: Bearer $TOKEN"

# 5. Access another protected endpoint
curl http://localhost:8080/api/v1/example \
  -H "Authorization: Bearer $TOKEN"

# 6. Logout
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer $TOKEN"
```

## Using API Key

```bash
# Get API key from database (default: default-api-key-12345)
API_KEY="default-api-key-12345"

# Access protected endpoint with API key
curl http://localhost:8080/api/v1/example \
  -H "X-API-Key: $API_KEY"

# Get user info with API key (note: API keys don't have associated users)
curl http://localhost:8080/api/v1/me \
  -H "X-API-Key: $API_KEY"
```

## Error Handling Examples

### Invalid Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "invalid", "password": "wrong"}'

# Response:
{
  "error": "Unauthorized",
  "message": "Invalid username or password"
}
```

### Missing Authentication

```bash
curl http://localhost:8080/api/v1/example

# Response:
{
  "error": "Unauthorized",
  "message": "Authentication required"
}
```

### Invalid Token

```bash
curl http://localhost:8080/api/v1/example \
  -H "Authorization: Bearer invalid-token"

# Response:
{
  "error": "Unauthorized",
  "message": "Authentication required"
}
```

### Not Found

```bash
curl http://localhost:8080/nonexistent

# Response:
{
  "error": "Not found",
  "path": "/nonexistent"
}
```

## Using with Different Tools

### Python Example

```python
import requests
import json

# Login
login_url = "http://localhost:8080/api/v1/auth/login"
login_data = {"username": "admin", "password": "admin123"}

response = requests.post(login_url, json=login_data)
token = response.json()["token"]

# Access protected endpoint
headers = {"Authorization": f"Bearer {token}"}
response = requests.get("http://localhost:8080/api/v1/me", headers=headers)
print(response.json())

# Logout
requests.post("http://localhost:8080/api/v1/auth/logout", headers=headers)
```

### JavaScript (Node.js) Example

```javascript
const axios = require('axios');

async function main() {
    try {
        // Login
        const loginResponse = await axios.post('http://localhost:8080/api/v1/auth/login', {
            username: 'admin',
            password: 'admin123'
        });
        
        const token = loginResponse.data.token;
        console.log('Token:', token);
        
        // Access protected endpoint
        const config = {
            headers: { Authorization: `Bearer ${token}` }
        };
        
        const meResponse = await axios.get('http://localhost:8080/api/v1/me', config);
        console.log('User info:', meResponse.data);
        
        // Logout
        await axios.post('http://localhost:8080/api/v1/auth/logout', {}, config);
        console.log('Logged out successfully');
    } catch (error) {
        console.error('Error:', error.response ? error.response.data : error.message);
    }
}

main();
```

### Postman Collection

You can import the following Postman collection to test the API:

```json
{
  "info": {
    "_postman_id": "home-ctrl-api",
    "name": "Home Control API",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Health Check",
      "request": {
        "method": "GET",
        "url": "http://localhost:8080/health"
      }
    },
    {
      "name": "Get Version",
      "request": {
        "method": "GET",
        "url": "http://localhost:8080/api/v1/version"
      }
    },
    {
      "name": "Login",
      "request": {
        "method": "POST",
        "url": "http://localhost:8080/api/v1/auth/login",
        "body": {
          "mode": "raw",
          "raw": "{\n  \"username\": \"admin\",\n  \"password\": \"admin123\"\n}"
        },
        "header": [
          {
            "key": "Content-Type",
            "value": "application/json"
          }
        ]
      }
    },
    {
      "name": "Get User Info",
      "request": {
        "method": "GET",
        "url": "http://localhost:8080/api/v1/me",
        "header": [
          {
            "key": "Authorization",
            "value": "Bearer {{token}}"
          }
        ]
      }
    },
    {
      "name": "Example Endpoint",
      "request": {
        "method": "GET",
        "url": "http://localhost:8080/api/v1/example",
        "header": [
          {
            "key": "Authorization",
            "value": "Bearer {{token}}"
          }
        ]
      }
    },
    {
      "name": "Logout",
      "request": {
        "method": "POST",
        "url": "http://localhost:8080/api/v1/auth/logout",
        "header": [
          {
            "key": "Authorization",
            "value": "Bearer {{token}}"
          }
        ]
      }
    }
  ]
}
```

## Configuration Examples

### Custom Configuration

Create a `config.yaml` file:

```yaml
server:
  host: "0.0.0.0"
  port: 8080

auth:
  users:
    admin: "admin123"
    user: "user123"
    guest: "guest123"
  session_ttl_hours: 8  # 8 hours session timeout

data_dir: "custom_data"
```

Then run the server with custom config:

```bash
HOME_CTRL_CONFIG=config.yaml go run cmd/home-ctrl/main.go
```

### Using Environment Variables

```bash
# Set custom data directory
export HOME_CTRL_DATA_DIR="/path/to/data"

# Run server
go run cmd/home-ctrl/main.go
```
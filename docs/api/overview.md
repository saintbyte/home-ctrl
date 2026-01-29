# API Overview

## Introduction

Home Control API provides a RESTful interface for managing home automation systems. The API is designed to be simple, secure, and extensible.

## Base URL

```
http://localhost:8080
```

## API Versioning

The API uses URL versioning with the `/api/v1/` prefix. This allows for backward compatibility when introducing new API versions.

## Authentication

The API supports two authentication methods:

1. **Bearer Token Authentication**: Obtain a token by logging in with username and password
2. **API Key Authentication**: Use API keys stored in the SQLite database

## Response Format

All responses are in JSON format with appropriate HTTP status codes.

## Error Handling

The API returns standard HTTP status codes and JSON error responses:

- `400 Bad Request`: Invalid request format
- `401 Unauthorized`: Authentication required or failed
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server-side error

## Rate Limiting

Currently, the API does not implement rate limiting, but this can be added in future versions.

## CORS

CORS is not explicitly configured but can be added using Gin middleware if needed for web applications.

## Security

- All sensitive endpoints require authentication
- Use HTTPS in production
- Rotate API keys regularly
- Set appropriate session TTL values
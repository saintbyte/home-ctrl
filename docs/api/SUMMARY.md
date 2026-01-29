# API Documentation Summary

## Overview

This directory contains comprehensive documentation for the Home Control API.

## Documentation Structure

### [Overview](overview.md)
- Introduction to the API
- Base URL and versioning
- Authentication methods
- Response formats
- Error handling

### [Authentication](authentication.md)
- Bearer Token Authentication
  - Login process
  - Using tokens
  - Logout
- API Key Authentication
  - Using API keys
  - Managing API keys
- Session management
- Configuration
- Security best practices

### [Endpoints](endpoints.md)
- Public endpoints (no authentication)
  - Health check
  - API version
- Authentication endpoints
  - Login
  - Logout
- Protected endpoints
  - User info
  - Example endpoint

### [Examples](examples.md)
- Complete workflow example
- Using API keys
- Error handling examples
- Python example
- JavaScript (Node.js) example
- Postman collection
- Configuration examples

### [Configuration](configuration.md)
- Configuration file format
- Server configuration
- Authentication configuration
- Data directory
- Using configuration files
- Environment variables
- Configuration best practices
- Advanced configuration

## Quick Start

1. **Check API health:**
   ```bash
   curl http://localhost:8080/health
   ```

2. **Login to get a token:**
   ```bash
   curl -X POST http://localhost:8080/api/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{"username": "admin", "password": "admin123"}'
   ```

3. **Access protected endpoints:**
   ```bash
   curl http://localhost:8080/api/v1/me \
     -H "Authorization: Bearer YOUR_TOKEN"
   ```

## API Features

- ✅ RESTful design
- ✅ JSON request/response
- ✅ Bearer token authentication
- ✅ API key authentication
- ✅ Session management
- ✅ Versioned endpoints
- ✅ Comprehensive error handling
- ✅ Configuration via YAML
- ✅ SQLite database backend

## Support

For issues or questions, please refer to the main [README.md](../../README.md) file or open an issue in the project repository.
# AI Agent Guide for Home Control Project

## Project Overview

This is a **Home Control System** built with Go, providing a RESTful API for home automation management. The project uses modern Go practices with a clean architecture.

## Project Structure

```
home-ctrl/
├── cmd/              # Entry points
│   └── home-ctrl/     # Main application
├── internal/         # Internal application logic
│   ├── app/          # Core application
│   ├── auth/         # Authentication
│   ├── config/       # Configuration
│   ├── database/     # SQLite database
│   └── server/       # HTTP server
│       └── v1/       # API v1 handlers
├── pkg/              # Reusable packages
│   └── utils/        # Utility functions
├── docs/             # Documentation
│   └── api/          # API documentation
├── examples/         # Example configurations
│   ├── config.yaml   # Configuration example
│   └── system/       # System service examples
├── Makefile          # Build automation
├── go.mod            # Go module
└── README.md         # Project documentation
```

## Key Components

### 1. HTTP Server (Gin Framework)
- **Port**: 8080 (configurable)
- **Base Path**: `/api/v1/`
- **Features**: 
  - Versioned API endpoints
  - Bearer token authentication
  - API key authentication
  - Middleware support

### 2. Authentication System
- **Methods**: 
  - Username/Password (from config)
  - API Keys (from SQLite database)
- **Session Management**: SQLite-based with TTL
- **Default Users**:
  - `admin:admin123`
  - `user:user123`

### 3. Database (SQLite)
- **Location**: `data/home-ctrl.db`
- **Tables**:
  - `api_keys` - API key storage
  - `sessions` - User sessions
- **Default API Key**: `default-api-key-12345`

### 4. Configuration
- **Format**: YAML
- **Default**: `config.yaml`
- **Environment Variable**: `HOME_CTRL_CONFIG`

## API Endpoints

### Public (No Auth)
- `GET /health` - Health check
- `GET /api/v1/version` - API version info

### Authentication
- `POST /api/v1/auth/login` - Get Bearer token
- `POST /api/v1/auth/logout` - Invalidate session

### Protected (Require Auth)
- `GET /api/v1/me` - Current user info
- `GET /api/v1/example` - Example protected endpoint

## Build & Run

### Build
```bash
# Build all versions
make build

# Build specific version
make build-linux-386    # 32-bit
make build-linux-amd64  # 64-bit
```

### Run
```bash
# Run with default config
go run cmd/home-ctrl/main.go

# Run with custom config
HOME_CTRL_CONFIG=custom.yaml go run cmd/home-ctrl/main.go
```

### Test
```bash
# Run all tests
make test

# Or directly
go test ./...
```

## Common Tasks

### Add New API Endpoint
1. Create handler in `internal/server/v1/`
2. Implement handler methods
3. Add to router in `internal/server/v1/router.go`

### Add New Database Table
1. Add migration in `internal/database/`
2. Create model struct
3. Add CRUD methods

### Add Configuration Option
1. Add to `internal/config/config.go`
2. Add to default config
3. Update documentation

## Architecture Principles

1. **Clean Architecture**: Separation of concerns
2. **Dependency Injection**: Components are loosely coupled
3. **Testability**: Easy to test individual components
4. **Extensibility**: Designed for future growth
5. **Versioning**: API versioning from the start

## Security Considerations

1. **Authentication**: Always required for sensitive endpoints
2. **HTTPS**: Recommended for production
3. **Input Validation**: All inputs are validated
4. **Error Handling**: Proper error responses
5. **Session Management**: Automatic cleanup of expired sessions

## Development Guidelines

1. **Code Style**: Follow Go conventions
2. **Testing**: Write tests for new functionality
3. **Documentation**: Update docs when adding features
4. **Commits**: Small, focused commits with clear messages
5. **Branching**: Use feature branches for development

## Troubleshooting

### Common Issues

1. **Port already in use**:
   ```bash
   lsof -i :8080
   kill <PID>
   ```

2. **Database permission issues**:
   ```bash
   mkdir -p data
   chmod 755 data
   ```

3. **Configuration errors**:
   ```bash
   # Check config syntax
   go run cmd/home-ctrl/main.go
   ```

## Future Enhancements

- [ ] Add more API endpoints for devices
- [ ] Implement device management
- [ ] Add MQTT integration
- [ ] Add WebSocket support
- [ ] Add OpenAPI/Swagger documentation
- [ ] Add rate limiting
- [ ] Add request logging
- [ ] Add health metrics

## For AI Agents

When working with this project:

1. **Understand the structure** before making changes
2. **Follow existing patterns** for consistency
3. **Write tests** for new functionality
4. **Update documentation** when adding features
5. **Keep changes minimal** and focused
6. **Use the Makefile** for common tasks
7. **Check the examples** for configuration

## Useful Commands

```bash
# Check Go version
go version

# Format code
go fmt ./...

# Check for issues
go vet ./...

# View dependencies
go mod graph

# Clean build artifacts
make clean
```

This guide provides a comprehensive overview for AI agents and developers working on the Home Control project. The structure is designed to be modular and extensible, making it easy to add new features while maintaining code quality.
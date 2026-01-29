# home-ctrl

A home automation control system.

## Quick Start

```bash
# Build the project
make build

# Run the application
go run cmd/home-ctrl/main.go
```

## System Service Setup

Example configuration files for running home-ctrl as a system service are available in the [examples/system](examples/system) directory:

- **systemd**: `examples/system/home-ctrl.service`
- **SystemV init**: `examples/system/home-ctrl.init`

See [examples/system/README.md](examples/system/README.md) for detailed installation instructions.

## Project Structure

```
home-ctrl/
├── cmd/              # Main application entry points
│   └── home-ctrl/     # Main application
├── internal/         # Internal application logic
│   └── app/          # Core application logic
├── pkg/              # Reusable packages
│   └── utils/        # Utility functions
├── examples/         # Example configurations
│   └── system/       # System service examples
├── bin/              # Build output (ignored by git)
├── Makefile          # Build automation
└── go.mod            # Go module definition
```

## Building

```bash
# Build for all platforms
make build

# Build for specific platform
make build-linux-386    # 32-bit Linux
make build-linux-amd64  # 64-bit Linux

# Clean build artifacts
make clean
```

## Testing

```bash
# Run all tests
make test

# Or directly with go
go test ./...
```

## API Documentation

Comprehensive API documentation is available in the [docs/api](docs/api) directory:

- [Overview](docs/api/overview.md) - Introduction and general information
- [Authentication](docs/api/authentication.md) - Authentication methods and usage
- [Endpoints](docs/api/endpoints.md) - Detailed endpoint documentation
- [Examples](docs/api/examples.md) - Practical usage examples
- [Configuration](docs/api/configuration.md) - Configuration options and best practices
- [Summary](docs/api/SUMMARY.md) - Complete documentation overview

## License

MIT License - see [LICENSE](LICENSE) for details.
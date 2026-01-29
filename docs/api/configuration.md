# Configuration

## Configuration File

The Home Control API can be configured using a YAML configuration file. By default, the application looks for `config.yaml` in the current directory, but you can specify a custom path using the `HOME_CTRL_CONFIG` environment variable.

## Configuration Options

### Server Configuration

```yaml
server:
  # Host to listen on (default: 127.0.0.1)
  host: "0.0.0.0"
  
  # Port to listen on (default: 8080)
  port: 8080
```

### Authentication Configuration

```yaml
auth:
  # User credentials (username: password)
  users:
    admin: "admin123"
    user: "user123"
  
  # Session TTL in hours (default: 24)
  session_ttl_hours: 24
```

### Data Directory

```yaml
# Data directory for SQLite database (default: "data")
data_dir: "data"
```

## Complete Example Configuration

```yaml
# Home Control Configuration
# Example configuration file for home-ctrl

# Server configuration
server:
  # Host to listen on (default: 127.0.0.1)
  host: "0.0.0.0"
  
  # Port to listen on (default: 8080)
  port: 8080

# Authentication configuration
auth:
  # User credentials (username: password)
  users:
    admin: "admin123"
    user: "user123"
    guest: "guest123"
  
  # Session TTL in hours (default: 24)
  session_ttl_hours: 8

# Data directory for SQLite database
data_dir: "data"
```

## Using Configuration File

### Default Configuration

The application will automatically look for `config.yaml` in the current directory:

```bash
go run cmd/home-ctrl/main.go
```

### Custom Configuration Path

Use the `HOME_CTRL_CONFIG` environment variable to specify a custom configuration file:

```bash
HOME_CTRL_CONFIG=/path/to/custom-config.yaml go run cmd/home-ctrl/main.go
```

### Command Line Example

```bash
# Create a custom configuration file
cat > custom-config.yaml <<EOF
server:
  host: "0.0.0.0"
  port: 8081

auth:
  users:
    admin: "securepassword"
    user: "userpassword"
  session_ttl_hours: 12

data_dir: "/var/lib/home-ctrl"
EOF

# Run with custom configuration
HOME_CTRL_CONFIG=custom-config.yaml go run cmd/home-ctrl/main.go
```

## Environment Variables

The application supports the following environment variables:

- `HOME_CTRL_CONFIG`: Path to the configuration file (default: `config.yaml`)
- `HOME_CTRL_DATA_DIR`: Override the data directory (optional)

## Configuration Best Practices

### Security

1. **Use Strong Passwords**: Always use strong, unique passwords for user accounts
2. **Limit Session TTL**: Set appropriate session expiration times based on your security requirements
3. **Secure Data Directory**: Ensure the data directory has proper permissions
4. **HTTPS**: Always use HTTPS in production environments

### Performance

1. **Port Selection**: Use standard ports (80, 443) or well-known ports for better compatibility
2. **Host Binding**: Bind to specific interfaces rather than `0.0.0.0` when possible

### Maintenance

1. **Backup Configuration**: Regularly backup your configuration files
2. **Version Control**: Store configuration files in version control
3. **Document Changes**: Document any configuration changes for future reference

## Configuration Validation

The application validates the configuration on startup and will fail with an error message if:

- The configuration file is malformed
- Required fields are missing
- Invalid values are provided

## Advanced Configuration

### Multiple Environments

You can maintain separate configuration files for different environments:

```bash
# Development
export HOME_CTRL_CONFIG=config.dev.yaml
go run cmd/home-ctrl/main.go

# Production
export HOME_CTRL_CONFIG=config.prod.yaml
go run cmd/home-ctrl/main.go
```

### Configuration Management

For production deployments, consider using configuration management tools like:

- Ansible
- Chef
- Puppet
- Kubernetes ConfigMaps
- Docker Secrets

### Dynamic Configuration

The current implementation loads configuration at startup. For dynamic configuration changes, you would need to:

1. Implement configuration reload endpoints
2. Use environment variable overrides
3. Implement configuration watchers
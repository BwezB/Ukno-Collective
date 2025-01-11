## Environment Variables

The services can be configured using environment variables. These override the default values and any values set in the configuration file.

### Common Variables
These variables affect the base configuration for any service:
- `ENVIRONMENT`: Controls environment-specific behaviors (`development` or `production`)
- `CONFIG_FILE_PATH`: Path to the configuration file (default: `config.yaml`)

### Server Configuration
Variables for configuring the gRPC and metrics servers:
- `SERVER_HOST`: Host address for the gRPC server
- `SERVER_PORT`: Port for the gRPC server
- `METRICS_HOST`: Host address for the Prometheus metrics server
- `METRICS_PORT`: Port for the metrics server
- `METRICS_PATH`: HTTP path for metrics endpoint

### Database Configuration
Variables for PostgreSQL database connection:
- `DB_HOST`: Database server hostname
- `DB_PORT`: Database server port
- `DB_USER`: Database username
- `DB_PASSWORD`: Database password [REQUIRED] - This field does not have a default value
- `DB_NAME`: Database name
- `DB_MAX_OPEN_CONNS`: Maximum number of open database connections
- `DB_MAX_IDLE_CONNS`: Maximum number of idle database connections
- `DB_CONN_MAX_LIFETIME`: Maximum lifetime of database connections (e.g., "5m", "1h")

### Logging Configuration
Variables for configuring the logging behavior:
- `LOGGER_ENVIRONMENT`: Logging preset (`development` or `production`)
- `LOG_LEVEL`: Minimum log level (`debug`, `info`, `warning`, or `error`)
- `LOG_ENCODING`: Log format (`json` or `console`)

### Health Check Configuration
- `HEALTH_CHECK_INTERVAL`: Interval between health checks (e.g., "5s", "1m")

### Auth Service Specific Variables
These variables are only used by the Auth service:
- `JWT_SECRET`: Secret key for signing JWT tokens [REQUIRED] - This field does not have a default value
- `JWT_EXPIRY`: JWT token expiration time (e.g., "24h", "168h")
- `AUTH_EMAIL`: Email identity for auth service
- `AUTH_PASSWORD`: Password for auth service [REQUIRED] - This field does not have a default value
- `GRAPH_HOST`: Host address of the graph service
- `GRAPH_PORT`: Port of the graph service

### Graph Service Specific Variables
These variables are only used by the Graph service:
- `AUTH_HOST`: Host address of the auth service
- `AUTH_PORT`: Port of the auth service

### Example Usage
```bash
# Basic setup
export ENVIRONMENT="development"
export SERVER_PORT="50051"
export DB_PASSWORD="your-secure-password"

# Enhanced logging for debugging
export LOG_LEVEL="debug"
export LOG_ENCODING="console"

# JWT configuration (Auth Service)
export JWT_SECRET="your-secure-jwt-secret"
export JWT_EXPIRY="24h"

# Service communication
export AUTH_HOST="localhost"
export AUTH_PORT="50051"
export GRAPH_HOST="localhost"
export GRAPH_PORT="50052"
```

### Command-Line Flags
All environment variables can also be set using command-line flags. The flags follow this pattern:
- Environment variable: `METRICS_HOST` → Flag: `--metrics-host`
- Environment variable: `DB_PASSWORD` → Flag: `--database-password`

For example:
```bash
./authservice --database-password="secure123" --log-level="debug"
./graphservice --server-port="50052" --metrics-host="localhost"
```

### Configuration Priority
The configuration values are applied in the following order (highest priority first):
1. Command-line flags
2. Environment variables
3. Configuration file
4. Default values

This means that environment variables will override values from the configuration file but can be overridden by command-line flags.
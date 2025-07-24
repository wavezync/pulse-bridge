# Pulse Bridge

Pulse Bridge is a lightweight, powerful uptime monitoring tool for your internal infrastructure (APIs, databases, etc.) and external platforms.

## How it works

Simply create a configuration file to define multiple services and databases to be checked at custom intervals. Pulse Bridge records the health status of each service and database, and provides a simple HTTP API to query their status.

#### Currently supports

## Monitoring

CLI > .env



## Configuration

The configuration file is a YAML file where you can define the services and databases you want to monitor.

Database monitors can be configured using a connection string or individual parameters (host, port, username, password, database name). The `driver` field is required to specify the database type (e.g., `postgres`, `mysql`, `mariadb`, `mssql`, `redis`).

You may also include a `query` field to run a custom SQL query for health checks, but it is not required.

Example configuration:

```yaml
monitors:
  # HTTP service monitoring
  - name: "HTTP Service"
    type: "http"
    interval: "30s"
    timeout: "5s"
    http:
      url: "http://helloworld-http:8080/ping"
      method: "GET"
      headers:
        Authorization: "Bearer secret-token"
        Content-Type: "application/json"

  # Postgres monitoring
  - name: "PostgreSQL Service"
    type: "database"
    interval: "30s"
    timeout: "10s"
    database:
      driver: "postgres"
      connection_string: "postgres://postgres:postgres@postgres-db:5432/monitoring?sslmode=disable"
      query: "SELECT 1"

  # MySQL monitoring
  - name: "MySQL Service"
    type: "database"
    interval: "30s"
    timeout: "10s"
    database:
      driver: "mysql"
      connection_string: "root:mysql@tcp(mysql-db:3306)/monitoring"
      query: "SELECT 1"

  # MariaDB monitoring
  - name: "MariaDB Service"
    type: "database"
    interval: "30s"
    timeout: "10s"
    database:
      driver: "mariadb"
      connection_string: "root:mariadb@tcp(mariadb-db:3306)/monitoring"
      query: "SELECT 1"

  # Redis monitoring
  - name: "Redis Service Primary"
    type: "database"
    interval: "5s"
    timeout: "5s"
    database:
      driver: "redis"
      database: "1"
      host: "redis-db"
      port: 6379
      password: "redispassword"

  # MSSQL monitoring
  - name: "MSSQL Service" 
    type: "database"
    interval: "30s"
    timeout: "10s"
    database:
      driver: "mssql"
      host: "mssql-db"
      port: 1433
      username: "SA"
      password: "Password1!"
      database: "master"
      query: "SELECT 1"
```

## Monitoring

You can check the status of your service from the pulse bridge API at the routes:


#### /monitor/services  
- List all monitored services

```json
[
  {
    "service": "HTTP Service",
    "status": "healthy",
    "type": "http",
    "last_check": "2025-07-24 11:56:01.918452021 +0000 UTC m=+0.357002662",
    "last_success": "2025-07-24 11:56:01.918443897 +0000 UTC m=+0.356994537",
    "metrics": {
      "response_time_ms": 81,
      "check_interval": "30s",
      "consecutive_successes": 1
    },
    "last_error": ""
  },
  {
    "service": "PostgreSQL Service",
    "status": "unhealthy",
    "type": "database",
    "last_check": "2025-07-24 11:56:01.891732112 +0000 UTC m=+0.330282750",
    "last_success": "",
    "metrics": {
      "response_time_ms": 50,
      "check_interval": "30s",
      "consecutive_successes": 0
    },
    "last_error": "failed to ping database: dial tcp 172.23.0.3:5432: connect: connection refused"
  }
 ]
```

#### /monitor/services/{monitor_name} 
- Get details of a specific service

```json
  {
    "service": "MariaDB Service",
    "status": "unhealthy",
    "type": "database",
    "last_check": "2025-07-24 11:56:01.89172233 +0000 UTC m=+0.330272963",
    "last_success": "",
    "metrics": {
      "response_time_ms": 33,
      "check_interval": "30s",
      "consecutive_successes": 0
    },
    "last_error": "failed to ping database: dial tcp 172.23.0.7:3306: connect: connection refused"
  },
```

## Contributing

We welcome contributions! If you have ideas, bug fixes, or improvements, please open an issue or submit a pull request.




Keep your systems transparent, your teams informed, and your users confident with Pulse Bridge – the heartbeat of your infrastructure. 🌊

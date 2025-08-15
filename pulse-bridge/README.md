# Pulse Bridge Helm Chart

A Helm chart for deploying Pulse Bridge - a comprehensive monitoring solution for HTTP services and databases.

## Overview

Pulse Bridge is a Go-based monitoring application that provides health checking capabilities for various services including HTTP endpoints, PostgreSQL, MySQL, MariaDB, Redis, and MSSQL databases. This Helm chart simplifies the deployment of Pulse Bridge on Kubernetes clusters.

## Prerequisites

- Kubernetes 1.16+
- Helm 3.0+

## Installation

### Quick Start

Install Pulse Bridge with your custom configuration:

```bash
helm install pulse-bridge ./pulse-bridge --set-file config.yaml=./my-config.yml
```

### Complete Installation with All Available Flags

```bash
helm install pulse-bridge ./pulse-bridge \
  --set-file config.yaml=./my-config.yml \
  --set replicaCount=2 \
  --set image.repository=pulse-bridge \
  --set image.tag=latest \
  --set image.pullPolicy=IfNotPresent \
  --set service.type=LoadBalancer \
  --set service.port=8080 \
```

## Configuration

### Configuration Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `replicaCount` | Number of Pulse Bridge replicas | `1` |
| `image.repository` | Docker image repository | `pulse-bridge` |
| `image.tag` | Docker image tag | `latest` |
| `image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `service.type` | Kubernetes service type | `LoadBalancer` |
| `service.port` | Service port | `8080` |
| `config.yaml` | Pulse Bridge monitoring configuration | See example below |
| `resources` | Pod resource requests and limits | `{}` |
| `nodeSelector` | Node labels for pod assignment | `{}` |
| `tolerations` | Toleration labels for pod assignment | `[]` |
| `affinity` | Affinity settings for pod assignment | `{}` |

### Configuration File Structure

Create a `my-config.yml` file with your monitoring targets:

```yaml
monitors:
  # HTTP service monitoring
  - name: "Web API Health Check"
    type: "http"
    interval: "30s"
    timeout: "5s"
    http:
      url: "https://api.example.com/health"
      method: "GET"
      headers:
        Content-Type: "application/json"
        Authorization: "Bearer your-token"

  # Database monitoring
  - name: "PostgreSQL Database"
    type: "database"
    interval: "30s"
    timeout: "10s"
    database:
      driver: "postgres"
      connection_string: "postgres://user:password@db-host:5432/dbname?sslmode=disable"
      query: "SELECT 1"

  # Redis monitoring
  - name: "Redis Cache"
    type: "database"
    interval: "15s"
    timeout: "5s"
    database:
      driver: "redis"
      host: "redis-host"
      port: 6379
      password: "redis-password"
      database: "0"
```

## API Endpoints

Once deployed, Pulse Bridge exposes the following endpoints:

- `GET /health` - Application health check
- `GET /monitor/services` - Status of all monitored services
- `GET /monitor/services/{service-name}` - Status of a specific service

## Upgrading

```bash
# Upgrade with new configuration
helm upgrade pulse-bridge ./pulse-bridge --set-file config.yaml=./updated-config.yml

# Upgrade with new image version
helm upgrade pulse-bridge ./pulse-bridge --set image.tag=v1.1.0
```

## Uninstalling

```bash
helm uninstall pulse-bridge
```

## Troubleshooting

### Check Pod Status

```bash
kubectl get pods -l app.kubernetes.io/name=pulse-bridge
```

### View Pod Logs

```bash
kubectl logs -l app.kubernetes.io/name=pulse-bridge -f
```

### Describe Pod for Events

```bash
kubectl describe pod -l app.kubernetes.io/name=pulse-bridge
```

### Check ConfigMap

```bash
kubectl get configmap pulse-bridge-config -o yaml
```

## Contributing

This Helm chart is part of the Pulse Bridge open-source project. Contributions are welcome!

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test the chart with different configurations
5. Submit a pull request

## Support

For issues and questions:

1. Check the main project repository: [wavezync/pulse-bridge](https://github.com/wavezync/pulse-bridge)
2. Review the application logs using `kubectl logs`
3. Open an issue on GitHub

## License

This chart is licensed under the same license as the main Pulse Bridge project.


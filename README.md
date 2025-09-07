# Air Quality Monitoring Server System

## Project Overview

This is an enterprise-grade air quality monitoring server system designed to receive, store, process, and analyze data from ESP32 air quality monitoring devices. The system adopts a simplified microservices architecture, developed in Go, with high availability, high performance, and easy deployment characteristics.

## System Features

### Core Functions
- **Data Reception**: Supports HTTP/WebSocket protocols for device data reception
- **Data Storage**: Uses MySQL for all data storage, Redis as cache and message queue
- **Real-time Processing**: Real-time data processing and alerts based on Redis Pub/Sub
- **Visualization**: Provides RESTful API interfaces for data query and visualization

### Technical Features
- **Simplified Architecture**: Removes complex time-series databases and message queues, uses lightweight MySQL+Redis solution
- **Monolithic Application**: Modular design, single deployment, easy maintenance
- **High Performance**: Supports high-concurrency data reception and processing
- **Easy Deployment**: One-click deployment using Docker Compose

## System Architecture

### Overall Architecture
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   ESP32 Device  │    │   Mobile APP    │    │   Web Admin     │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          │ HTTP/WebSocket       │ HTTP/WebSocket       │ HTTP/WebSocket
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                    ┌─────────────┴─────────────┐
                    │  Air Quality Server       │
                    │  (Monolithic Architecture)│
                    └─────────────┬─────────────┘
                                 │
        ┌────────────────────────┼────────────────────────┐
        │                        │                        │
┌───────▼────────┐    ┌─────────▼─────────┐    ┌─────────▼─────────┐
│  Message Queue │    │  Relational DB    │    │   Redis Cache     │
│ (Redis Pub/Sub)│    │   (MySQL)         │    │   (Redis)         │
└────────────────┘    └───────────────────┘    └───────────────────┘
```

## Technology Stack

### Backend Technologies
- **Language**: Go 1.21+
- **Framework**: Gin (HTTP) + gRPC
- **Database**: MySQL 8.0 + Redis 7.0
- **Message Queue**: Redis Pub/Sub
- **Monitoring**: Prometheus + Grafana
- **Logging**: Structured logging (JSON format)

### Deployment and Operations
- **Containerization**: Docker + Docker Compose
- **Reverse Proxy**: Nginx
- **CI/CD**: GitLab CI/CD
- **Monitoring**: Prometheus + Grafana

## Quick Start

### Environment Requirements

#### Production Environment
- **Docker**: 24.0+ (Recommended 28.4+)
- **Docker Compose**: V2 (Integrated in Docker) or V1 (Standalone installation)
- **Operating System**: Linux (Ubuntu 20.04+, CentOS 8+), macOS, Windows
- **Memory**: Minimum 2GB available memory
- **Storage**: Minimum 5GB available disk space

#### Development Environment
- **Go**: 1.21+
- **Docker**: 24.0+ (Recommended 28.4+)
- **Docker Compose**: V2 or V1
- **Git**: 2.0+

#### Port Requirements
- **3308**: MySQL Database (Production)
- **3307**: MySQL Database (Development)
- **6381**: Redis Cache (Production)
- **6380**: Redis Cache (Development)
- **8082**: Web Application (Production)
- **8083**: Web Application (Development)
- **1883**: MQTT Broker

> **Note**: If the host machine has ports 3306, 6379, 8080 occupied, the system will automatically use the above alternative ports.

### Deployment Steps

1. **Clone Project**
```bash
git clone <repository-url>
cd air-quality-server
```

2. **Check Docker Environment**
```bash
# Check Docker version
docker --version

# Check Docker Compose version
docker compose version  # V2 syntax
# or
docker-compose --version  # V1 syntax
```

3. **Start Services**

**Using Docker Compose V2 (Recommended):**
```bash
# Start all services
docker compose up --build -d
```

**Using Docker Compose V1:**
```bash
# Start all services
docker-compose up --build -d
```

**Using Smart Start Script:**
```bash
# Automatically select available Docker Compose version
chmod +x scripts/docker/start-services.sh
./scripts/docker/start-services.sh
```

4. **Verify Deployment**
```bash
# Check service status (V2)
docker compose ps

# Check service status (V1)
docker-compose ps

# Test API
curl http://localhost:8082/health
```

5. **Access Services**
- **Web Interface**: http://localhost:8082
- **MySQL**: localhost:3308
- **Redis**: localhost:6381
- **MQTT**: localhost:1883

### Development Environment

1. **Install Dependencies**
```bash
go mod download
```

2. **Configure Environment**
```bash
# Use default configuration file
export CONFIG_FILE=config/config.yaml

# Or configure using environment variables
export DB_HOST=localhost
export DB_PORT=3306
export REDIS_HOST=localhost
export REDIS_PORT=6379
```

3. **Start Development Environment**
```bash
make dev
```

4. **Run Tests**
```bash
make test
```

## Project Structure

```
air-quality-server/
├── cmd/                # Application entry points
│   └── server/
│       └── main.go     # Main program entry
├── internal/           # Internal packages
│   ├── models/        # Data models
│   ├── services/      # Business logic layer
│   ├── repositories/  # Data access layer
│   ├── handlers/      # HTTP handlers
│   ├── middleware/    # Middleware
│   ├── config/        # Configuration management
│   └── utils/         # Utility packages
├── config/            # Configuration files
├── scripts/           # Script files
├── docs/              # Documentation
├── docker-compose.yml # Docker orchestration file
├── Dockerfile         # Docker build file
├── Makefile           # Build scripts
└── go.mod             # Go module file
```

## Documentation

### System Documentation
- [System Design Document](docs/system_design.md) - Detailed system architecture design
- [Module Interface Document](docs/module_interfaces.md) - Module interface definitions
- [Database Design Document](docs/database_design.md) - Database table structure design
- [Deployment Guide](docs/deployment_guide.md) - Detailed deployment and operations guide

### Docker Related Documentation
- [Docker Deployment Guide](docs/docker_guide.md) - Complete Docker deployment guide

## Troubleshooting

### Docker Related Issues

#### 1. ContainerConfig Error
```bash
# Use smart start script
chmod +x scripts/docker/start-services.sh
./scripts/docker/start-services.sh
```

#### 2. Port Conflicts
```bash
# Check port usage
lsof -i :3308
lsof -i :6381
lsof -i :8082

# Use smart start script
chmod +x scripts/docker/start-services.sh
./scripts/docker/start-services.sh
```

#### 3. Go Module Download Timeout (China)
```bash
# Set Go proxy environment variables
export GOPROXY=https://goproxy.cn,direct
export GOSUMDB=sum.golang.google.cn

# Then start services
docker compose up --build -d
```

#### 4. Docker Version Too Old
```bash
# Refer to Docker upgrade guide
# docs/docker_upgrade_ubuntu24.md
```

#### 5. docker-compose Command Not Found
```bash
# Use Docker Compose V2
docker compose up --build -d

# Or use smart start script
chmod +x scripts/docker/start-services.sh
./scripts/docker/start-services.sh
```

#### 6. Docker Compose Version Warning
```bash
# Warning: the attribute `version` is obsolete
# Solution: Removed version field from all docker-compose files
# Verify fix:
docker compose config
```

### Service Management

#### Start Services
```bash
# Production environment
docker compose up --build -d

# Development environment
docker compose -f docker-compose.dev.yml up --build -d
```

#### Stop Services
```bash
# Stop all services
docker compose down

# Stop and remove data volumes
docker compose down -v
```

#### View Logs
```bash
# View all service logs
docker compose logs -f

# View specific service logs
docker compose logs -f air-quality-server
```

## Contributing

1. Fork the project
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contact

For questions or suggestions, please contact us through:

- Project Issues: [GitHub Issues](https://github.com/coolham/air-quality-server/issues)
- Email: xingshizhai@gmail.com

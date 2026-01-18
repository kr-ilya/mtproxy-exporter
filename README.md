# MTProxy Exporter for Prometheus

A modern, efficient Prometheus exporter for MTProxy (Telegram's MTProto Proxy) written in Go 1.25+.

## Features

- 🚀 **Modern Architecture**: Clean, maintainable code following Go best practices
- 📊 **Comprehensive Metrics**: Exports all important MTProxy statistics to Prometheus
- 🔧 **Configurable**: Support for environment variables and command-line flags
- 🧪 **Well-Tested**: Extensive unit tests with high coverage
- 🐳 **Docker Ready**: Multi-stage Dockerfile for minimal image size
- 🔄 **Graceful Shutdown**: Proper signal handling and cleanup
- 🏥 **Health Checks**: Built-in health and readiness endpoints

## Quick Start

## Configuration

Configuration can be provided via environment variables or command-line flags:

| Flag | Environment Variable | Default | Description |
|------|---------------------|---------|-------------|
| `--listen.address` | `LISTEN_ADDRESS` | `:9330` | Address to listen on |
| `--mtproxy.url` | `MTPROXY_URL` | `http://localhost:8888` | MTProxy stats URL |
| `--mtproxy.timeout` | `MTPROXY_TIMEOUT` | `10s` | HTTP client timeout |
| `--log.level` | `LOG_LEVEL` | `info` | Log level (debug, info, warn, error) |

### Using Docker

```bash
docker run -d \
  -p 9330:9330 \
  -e MTPROXY_URL=http://your-mtproxy:8888 \
  mtproxy-exporter:latest
```

### Connecting to MTProxy Running on the Host

If your MTProxy is running directly on the host machine (not in Docker), the exporter container cannot access it via localhost by default. Use one of the following methods depending on your operating system.

#### Linux

Use Docker host network mode:
```
docker run -d \
  --name mtproxy-exporter \
  --network host \
  -e MTPROXY_URL=http://127.0.0.1:8888 \
  mtproxy-exporter:latest
```

#### MacOS / Windows

Use the special hostname `host.docker.internal`.


### Connecting to MTProxy in Docker

If your MTProxy runs in a Docker container (e.g., [mtproxy-docker](https://github.com/kr-ilya/mtproxy-docker)), you need to connect both containers to the same network:

#### 1️⃣ Create network (if not exists)

```bash
docker network create mtproxy-net
```

#### 2️⃣ Connect MTProxy to the network

If MTProxy is already running, connect it without rebuild:

```bash
docker network connect mtproxy-net mtproxy
```

> **Note:** Replace `mtproxy` with your actual container name (check with `docker ps`)

Or use Docker Compose with the `networks` section (specify the network as `external: true`).

#### 3️⃣ Run exporter with the same network

```bash
docker run -d \
  --name mtproxy-exporter \
  --network mtproxy-net \
  -p 9330:9330 \
  -e MTPROXY_URL=http://mtproxy:8888 \
  mtproxy-exporter
```

> **Key point:** `mtproxy` is the hostname of the container, Docker will resolve it to the container's IP automatically

### Using Binary

```bash
# Install
go install github.com/kr-ilya/mtproxy-exporter/cmd/mtproxy-exporter@latest

# Run
mtproxy-exporter --mtproxy.url=http://localhost:8888
```

## Exported Metrics

For a complete list of all exported metrics with descriptions, see [METRICS.md](METRICS.md).

## Prometheus Configuration

Add to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'mtproxy'
    static_configs:
      - targets: ['localhost:9330']
    scrape_interval: 15s
```

## Development

### Prerequisites

- Go 1.25+
- Make (optional)

### Building

```bash
# Build binary
make build

# Run tests
make test

# Run tests with coverage
make test-coverage

# Build Docker image
make docker-build

# Run locally
make run
```

## Testing

Run all tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

## License

MIT License

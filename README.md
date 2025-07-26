# HTTP Proxy

**HTTP Proxy** is a lightweight HTTP reverse proxy written in Go. It listens on a given address and port, and forwards all HTTP requests to a specified target server. It's ideal for simple debugging, local forwarding, and lightweight HTTP routing tasks.

## 📦 Install

Make sure you have [Go installed](https://golang.org/dl/), then:

```bash
go install github.com/sebps/http-proxy@latest
```

## 🚀 Features

- Forwards all incoming HTTP requests to a target host and port
- Logs request method, path, and body

## 🛠️ Usage

```bash
http-proxy [options]
```

### Options

| Flag              | Description                                      | Default     |
|-------------------|--------------------------------------------------|-------------|
| `--targetHost`    | Target host to forward requests to (required)    |             |
| `--targetPort`    | Target port to forward requests to (required)    |             |
| `--sourceAddr`    | Source address to bind the proxy server          | `localhost` |
| `--sourcePort`    | Source port to listen on                         | `80`        |
| `-h`, `--help`    | Show help message and exit                       |             |

## 🧪 Examples

Start a proxy that listens on port 8888 and forwards to `example.com:8080`:

```bash
./http-proxy --targetHost example.com --targetPort 8080 --sourcePort 8888
```

Bind to all interfaces instead of just `localhost`:

```bash
./http-proxy --targetHost example.com --targetPort 8080 --sourceAddr 0.0.0.0
```

## 🔎 Example Log Output

```
🚀 HTTP Proxy starting on http://localhost:8888 → http://example.com:8080
📥 GET request on path: /
🔎 Body:
```
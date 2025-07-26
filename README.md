# HTTP Proxy

**HTTP Proxy** is a lightweight HTTP reverse proxy written in Go. It listens on a given address and port, and forwards all HTTP requests to a specified target server. It's ideal for simple debugging, local forwarding, and lightweight HTTP routing tasks.

## 📦 Install

```bash
go install github.com/sebps/http-proxy
```

## 🚀 Features

- Forwards all incoming HTTP requests to a target host and port
- Logs request method, path, and body
- Optional automatic CORS headers via `--withCors`
- Clean CLI with `--help`, `--targetHost`, `--targetPort`, `--targetProtocol`, `--sourceAddr`, and `--sourcePort`
- Simple and minimal — no external dependencies beyond Go stdlib

## 🛠️ Usage

```bash
http-proxy [options]
```

### Options

| Flag                  | Description                                          | Default     |
|-----------------------|------------------------------------------------------|-------------|
| `--targetHost`        | Target host to forward requests to (required)        |             |
| `--targetPort`        | Target port to forward requests to                   | `80` or `443` based on protocol |
| `--targetProtocol`    | Protocol to use to reach the target (`http` or `https`) | `http`    |
| `--sourceAddr`        | Source address to bind the proxy server              | `localhost` |
| `--sourcePort`        | Source port to listen on                             | `80`        |
| `--withCors`          | Enable automatic CORS headers                        | `false`     |
| `-h`, `--help`        | Show help message and exit                           |             |

## 🧪 Examples

Start a proxy that listens on port 8888 and forwards to `https://example.com` with CORS support:

```bash
./http-proxy --targetHost example.com --targetProtocol https --sourcePort 8888 --withCors
```

Bind to all interfaces instead of just `localhost`:

```bash
./http-proxy --targetHost example.com --targetPort 443 --targetProtocol https --sourceAddr 0.0.0.0
```

## 🔎 Example Log Output

```
🚀 HTTP Proxy starting on http://localhost:8888
🔁 Forwarding all requests to https://example.com:443
📥 GET request on path: /api
🔎 Body: {"query": "hello"}
```
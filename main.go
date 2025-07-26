package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

// Proxy configuration
var (
	targetHost     string
	targetPort     string
	targetProtocol string
	sourceAddr     string
	sourcePort     string
	withCors       bool
)

func main() {
	// Custom help message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `HTTP Proxy

This tool starts an HTTP proxy server that forwards all incoming requests
to a specified target host and port using the specified protocol.

Usage:
  http-proxy [options]

Options:
  --targetHost <host>      Target host to proxy to (required)
  --targetPort <port>      Target port to proxy to (default: 80 or 443 based on protocol)
  --targetProtocol <proto> Target protocol (http or https) (default: http)
  --sourceAddr <addr>      Source address to bind to (default: localhost)
  --sourcePort <port>      Source port to listen on (default: 80)
  --withCors               Enable automatic CORS headers (default: false)
  -h, --help               Show this help message and exit
`)
	}

	// Flags
	flag.StringVar(&targetHost, "targetHost", "", "Target host to proxy to")
	flag.StringVar(&targetPort, "targetPort", "", "Target port to proxy to")
	flag.StringVar(&targetProtocol, "targetProtocol", "http", "Protocol to use for the target server (http or https)")
	flag.StringVar(&sourceAddr, "sourceAddr", "localhost", "Source address to bind to")
	flag.StringVar(&sourcePort, "sourcePort", "80", "Source port to listen on")
	flag.BoolVar(&withCors, "withCors", false, "Enable automatic CORS headers")
	help := flag.Bool("h", false, "Show help")
	flag.BoolVar(help, "help", false, "Show help (shorthand)")

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}
	if targetHost == "" {
		fmt.Fprintln(os.Stderr, "Error: --targetHost is required")
		flag.Usage()
		os.Exit(1)
	}
	if targetProtocol != "http" && targetProtocol != "https" {
		fmt.Fprintln(os.Stderr, "Error: --targetProtocol must be 'http' or 'https'")
		flag.Usage()
		os.Exit(1)
	}
	if targetPort == "" {
		if targetProtocol == "http" {
			targetPort = "80"
		} else {
			targetPort = "443"
		}
	}

	targetURL, err := url.Parse(fmt.Sprintf("%s://%s:%s", targetProtocol, targetHost, targetPort))
	if err != nil {
		log.Fatalf("❌ Invalid target URL: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
		req.Host = targetURL.Host

		if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
			req.Header.Set("X-Forwarded-For", clientIP)
		}
	}

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		proxy.ServeHTTP(w, r)
	})

	if withCors {
		handler = withCORS(handler)
	}

	http.Handle("/", handler)

	address := fmt.Sprintf("%s:%s", sourceAddr, sourcePort)
	log.Printf("🚀 HTTP Proxy starting on http://%s", address)
	log.Printf("🔁 Forwarding all requests to %s", targetURL)

	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func logRequest(r *http.Request) {
	log.Printf("📥 %s request on path: %s", r.Method, r.URL.Path)
	bodyBytes, _ := io.ReadAll(r.Body)
	r.Body.Close()
	r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	log.Printf("🔎 Body: %s", string(bodyBytes))
}

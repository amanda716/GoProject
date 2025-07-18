package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func hello(w http.ResponseWriter, r *http.Request) {
	// Log the incoming request details
	clientIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		clientIP = r.RemoteAddr
	}

	log.Printf("Request from %s - %s %s", clientIP, r.Method, r.URL.Path)
	fmt.Fprintf(w, "Hello")
}

// Custom logger that prints more connection info
type ConnectionLogger struct {
	handler http.Handler
}

func (l *ConnectionLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	l.handler.ServeHTTP(w, r)
	duration := time.Since(start)

	clientIP, clientPort, _ := net.SplitHostPort(r.RemoteAddr)
	log.Printf("Connection: src=%s sport=%s dst=%s dport=8080 method=%s path=%s duration=%v",
		clientIP, clientPort, r.Host, r.Method, r.URL.Path, duration)
}

func main() {
	// Set up logging
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Printf("Server starting on 0.0.0.0:8080")

	// Register route handler
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)

	// Wrap with connection logger
	loggedHandler := &ConnectionLogger{handler: mux}

	// Start server with custom handler
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", loggedHandler))
}

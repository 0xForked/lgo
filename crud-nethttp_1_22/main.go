package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/0xForked/crud-nethttp-1_22/todo"
)

var (
	port, env *string
)

const (
	local      = "local"
	staging    = "staging"
	production = "production"
)

func init() {
	port = flag.String("port", "8000", "server port")
	env = flag.String("env", "local",
		"server environment [local*, staging, production]")
	flag.Parse()
}

// applyMiddleware applies middleware to the handler
func applyMiddleware(handler http.Handler) http.Handler {
	if *env == local {
		handler = loggingMiddleware(handler)
	}
	return handler
}

// loggingResponseWriter is a wrapper around http.ResponseWriter that captures the status code
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code before calling the underlying WriteHeader method
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// loggingMiddleware is a simple middleware that logs incoming requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := &loggingResponseWriter{ResponseWriter: w}
		next.ServeHTTP(lrw, r)
		log.Printf("Request: %s %s | Status: %d\n",
			r.Method, r.URL.Path, lrw.statusCode)
	})
}

func main() {
	// init mux
	mux := http.NewServeMux()
	// handle todos routes
	todo.NewHandler(mux)
	// start server
	if err := http.ListenAndServe(
		fmt.Sprintf(":%s", *port),
		applyMiddleware(mux),
	); err != nil {
		log.Fatal("server failed to start: ", err)
	}
}

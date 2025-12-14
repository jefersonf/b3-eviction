package main

import (
	"b3e/internal/core/command"
	"b3e/internal/transport/rest"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/stretchr/testify/mock"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

var (
	DefaultListenAddr = "localhost:8080"
)

type MockBus struct {
	mock.Mock
}

func (b *MockBus) Publish(ctx context.Context, cmd command.CastVote) error {
	args := b.Called(ctx, cmd)
	return args.Error(0)
}

func main() {

	listenAddr := flag.String("listen-addr", DefaultListenAddr, fmt.Sprintf("Application Listen address (default is %s).", DefaultListenAddr))
	flag.Parse()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	router := http.NewServeMux()
	router.HandleFunc("/", WrapMiddleware(
		HealthCheckHandler(),
		AllowedResources("/"),
		Method(http.MethodGet),
		Logger(),
	))

	mockBus := new(MockBus) // Replace with actual bus implementation
	mockBus.On("Publish", mock.Anything, mock.AnythingOfType("command.CastVote")).Return(nil)

	voteHandler := rest.NewVoteHandler(mockBus)

	router.HandleFunc("/vote", WrapMiddleware(voteHandler.HandleVote, AllowedResources("/vote"), Method(http.MethodPost), Logger()))

	server := &http.Server{
		Addr:        *listenAddr,
		Handler:     router,
		ReadTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("Starting server on %s\n", *listenAddr)
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Could not listen on %s: %v\n", *listenAddr, err)
		}
	}()

	<-stopChan
	log.Println("Shutting down API server")

	ctx, shutdownServer := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownServer()
	server.Shutdown(ctx)
	log.Println("Server gracefully stopped")
}

func HealthCheckHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"status\": \"ok\"}"))
	}
}

func Logger() Middleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			defer func() {
				log.Printf("%s %s %v\n", r.Method, r.URL.Path, time.Since(start))
			}()
			h(w, r)
		}
	}
}

func Method(m string) Middleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != m {
				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
				return
			}
			h(w, r)
		}
	}
}

func AllowedResources(resources ...string) Middleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			for _, resource := range resources {
				if r.URL.Path == resource {
					h(w, r)
					return
				}
				notFound(w, r)
			}
		}
	}
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"message\": \"resource not found\"}"))
}

func WrapMiddleware(h http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

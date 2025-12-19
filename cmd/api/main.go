package main

import (
	"b3e/internal/bus"
	"b3e/internal/service/analytics"
	"b3e/internal/service/postgres"
	"b3e/internal/storage"
	"b3e/internal/transport/rest"
	"cmp"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	DefaultListenAddr = cmp.Or(os.Getenv("API_ADDR"), ":8080")
	DefaultQueueAddr  = cmp.Or(os.Getenv("QUEUE_ADDR"), "redis:6379")
	DefaultQueueName  = cmp.Or(os.Getenv("QUEUE_NAME"), "votes")
)

func main() {

	listenAddr := flag.String("listen-addr", DefaultListenAddr, fmt.Sprintf("Application listen address (default is %s).", DefaultListenAddr))
	queueAddr := flag.String("queue-addr", DefaultQueueAddr, fmt.Sprintf("Queue address (default is %s).", DefaultQueueAddr))
	queueName := flag.String("queue-name", DefaultQueueName, fmt.Sprintf("Queue name (default is %s).", DefaultQueueName))
	flag.Parse()

	// Redis queue instance
	voteStream := bus.NewStreamBus(*queueAddr, *queueName)
	// Postgres DB instance
	dbConnectionPool, err := postgres.NewConnectionPool()
	if err != nil {
		log.Fatalf("Failed to spin up the database connection pool: %v", err)
	}
	defer postgres.Shutdown(dbConnectionPool)

	voteRepo := storage.NewVoteRepository(dbConnectionPool)

	statsService := analytics.NewStatsService(voteStream)
	timelyStatsService := analytics.NewTimelyStatsService(voteRepo)

	voteHandler := rest.NewVoteHandler(voteStream)
	analyticsHandler := rest.NewVoteAnalyticsHandler(statsService, timelyStatsService)

	router := http.NewServeMux()
	router.HandleFunc("POST /vote", voteHandler.HandleVote)
	router.HandleFunc("GET /stats/{$}", rest.HandleHealthCheck)
	router.HandleFunc("GET /stats/{evictionId}", analyticsHandler.HandleEvictionStats)
	router.HandleFunc("GET /analytics/hourly", analyticsHandler.HandleHourlyStats)
	router.HandleFunc("GET /analytics/minutely", analyticsHandler.HandleMinutelylyStats)

	server := &http.Server{
		Addr:         *listenAddr,
		Handler:      rest.JSON(rest.CORS(router)),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("Starting voting API server on %s\n", *listenAddr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("could not listen on %s: %v\n", *listenAddr, err)
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	<-stopChan
	log.Println("Shutting down voting API server...")

	// A deadline to wait for active requests to complete
	ctx, shutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdown()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Voting API server forced to shutdown: %v\n", err)
	}

	log.Println("Server gracefully stopped")
}

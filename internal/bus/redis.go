package bus

import (
	"b3e/internal/core/command"
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

// StreamBus implements command.Publisher using Redis Streams.
type StreamBus struct {
	Client *redis.Client
	Stream string
}

// NewStreamBus creates a new Redis StreamBus.
func NewStreamBus(addr, streamName string) *StreamBus {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0, // use default DB
	})
	return &StreamBus{Client: rdb, Stream: streamName}
}

// Publish sends cast vote command to Redis stream.
func (r *StreamBus) Publish(ctx context.Context, cmd command.CastVote) error {
	args := &redis.XAddArgs{
		Stream: r.Stream,
		ID:     "*", // auto-generated ID
		Values: map[string]interface{}{
			"nominee_id":  cmd.NomineeID,
			"eviction_id": cmd.EvictionID,
			"timestamp":   cmd.Timestamp.Unix(),
		},
	}

	id, err := r.Client.XAdd(ctx, args).Result()
	if err != nil {
		return fmt.Errorf("failed to publish to stream: %w", err)
	}

	log.Printf("Published vote to %s as stream ID %s\n", cmd.NomineeID, id)
	return nil
}

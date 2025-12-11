package redisadapter

import (
	"b3e/internal/core/command"
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type RedisBus struct {
	client *redis.Client
	queue  string
}

func NewRedisBus(addr, queueName string) *RedisBus {
	rdb := redis.NewClient(&redis.Options{Addr: addr})
	return &RedisBus{client: rdb, queue: queueName}
}

// Publish (used by Ingestor)
func (r *RedisBus) Publish(ctx context.Context, cmd command.CastVote) error {
	data, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return r.client.LPush(ctx, r.queue, data).Err()
}

// Subscribe (used by Processor)
// this converts the Redis polling loop into a Go channel.
func (r *RedisBus) Subscribe(ctx context.Context) (<-chan command.CastVote, error) {
	ch := make(chan command.CastVote)

	go func() {
		defer close(ch)
		for {
			// BRPPop blocks until an item is available or timeout
			// 0 means block indefinitely
			result, err := r.client.BRPop(ctx, 0, r.queue).Result()
			if err != nil {
				// If context is canceled (shutdown), stop via return
				if ctx.Err() != nil {
					return
				}
				continue // retry on temp redis error
			}

			// result[0] is key, result[1] is value
			var cmd command.CastVote
			if err := json.Unmarshal([]byte(result[1]), &cmd); err != nil {
				select {
				case ch <- cmd:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return ch, nil
}

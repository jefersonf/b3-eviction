package analytics

import (
	"b3e/internal/bus"
	"b3e/internal/domain"
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
)

type redisStatsService struct {
	rdb       *redis.Client
	queueName string
}

func NewStatsService(streamBus *bus.StreamBus) StatsProvider {
	return &redisStatsService{rdb: streamBus.Client, queueName: streamBus.Stream}
}

// EvictionStats fetches counts for all nominees in a specific eviction and total vote count.
func (v *redisStatsService) EvictionStats(ctx context.Context, evictionID string) (domain.EvictionStats, error) {

	// Pattern: {queueName}:{evictionID}:*
	pattern := fmt.Sprintf("%s:%s:*", v.queueName, evictionID)
	iter := v.rdb.Scan(ctx, 0, pattern, 0).Iterator()

	evictionStats := domain.EvictionStats{
		VotesByNominee: make([]domain.VoteCount, 0),
	}

	// Collect all matching keys first
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return evictionStats, err
	}

	// Optimization: If no keys found, return empty early
	if len(keys) == 0 {
		return evictionStats, nil
	}

	// Optimization: Using MGET (Multi-Get) to fetch values in ONE roud-trip
	values, err := v.rdb.MGet(ctx, keys...).Result()
	if err != nil {
		return evictionStats, err
	}

	evictionStats.Evicted = domain.VoteCount{NomineeID: "unknown"}

	// Map keys and values to build eviction final stats
	for i, val := range values {
		if val == nil {
			continue
		}
		count, _ := strconv.ParseInt(val.(string), 10, 64)
		// Extract strictly the nominee_id from "votes:eviction_1:nominee_A"
		// Split by ":" and take the last part
		parts := strings.Split(keys[i], ":")
		if len(parts) >= 3 {
			evictionStats.VotesByNominee = append(
				evictionStats.VotesByNominee,
				domain.VoteCount{NomineeID: parts[2], Votes: count})

			evictionStats.TotalVotes += count
			if count > evictionStats.Evicted.Votes {
				evictionStats.Evicted.Votes = count
				evictionStats.Evicted.NomineeID = parts[2]
			}
		}
	}

	return evictionStats, nil
}

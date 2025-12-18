package analytics

import (
	"b3e/internal/domain"
	"context"
)

type StatsProvider interface {
	EvictionStats(ctx context.Context, evictionID string) (domain.EvictionStats, error)
}

type TimelyStatsProvider interface {
	HourlyStats(ctx context.Context) ([]domain.TimelyStat, error)
	MinutelyStats(ctx context.Context) ([]domain.TimelyStat, error)
}

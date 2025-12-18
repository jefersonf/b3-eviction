package analytics

import (
	"b3e/internal/domain"
	"b3e/internal/storage"
	"context"
)

type postgresService struct {
	repo storage.VoteRepository
}

func NewTimelyStatsService(repository storage.VoteRepository) TimelyStatsProvider {
	return &postgresService{repo: repository}
}

// HourlyStats ...
func (s *postgresService) HourlyStats(ctx context.Context) ([]domain.TimelyStat, error) {
	return s.repo.GetHourlyStats()
}

// MinutelyStats ...
func (s *postgresService) MinutelyStats(ctx context.Context) ([]domain.TimelyStat, error) {
	return s.repo.GetMinutelyStats()
}

package storage

import "b3e/internal/domain"

type VoteRepository interface {
	GetHourlyStats() ([]domain.TimelyStat, error)
	GetMinutelyStats() ([]domain.TimelyStat, error)
}

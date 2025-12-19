package storage

import (
	"b3e/internal/domain"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type voteRepository struct {
	pool *pgxpool.Pool
}

// NewVoteRepository accepts a connection pool, not a single connection.
func NewVoteRepository(pool *pgxpool.Pool) VoteRepository {
	return &voteRepository{pool: pool}
}

func (r *voteRepository) GetHourlyStats() ([]domain.TimelyStat, error) {
	query := `
        SELECT bucket_hour, nominee_id, total_votes 
        FROM votes_hourly 
        WHERE bucket_hour > NOW() - INTERVAL '24 hours'
        ORDER BY bucket_hour ASC, nominee_id ASC;
    `
	rows, err := r.pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]domain.TimelyStat, 0)
	for rows.Next() {
		var s domain.TimelyStat
		if err := rows.Scan(&s.Timedate, &s.NomineeID, &s.Votes); err != nil {
			continue
		}
		stats = append(stats, s)
	}
	return stats, nil
}

func (r *voteRepository) GetMinutelyStats() ([]domain.TimelyStat, error) {
	query := `
        SELECT bucket_minute, nominee_id, votes 
        FROM votes_minutely 
        WHERE bucket_minute > NOW() - INTERVAL '24 hours'
        ORDER BY bucket_minute ASC, nominee_id ASC;
    `
	rows, err := r.pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]domain.TimelyStat, 0)
	for rows.Next() {
		var s domain.TimelyStat
		if err := rows.Scan(&s.Timedate, &s.NomineeID, &s.Votes); err != nil {
			continue
		}
		stats = append(stats, s)
	}
	return stats, nil
}

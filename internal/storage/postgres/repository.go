package postgres

import (
	"b3e/internal/domain"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type VoteRepository struct {
	pool *pgxpool.Pool
}

// NewVoteRepository accepts a connection pool, not a single connection.
// This is critical for high-throughput concurrency.
func NewVoteRepository(pool *pgxpool.Pool) *VoteRepository {
	return &VoteRepository{pool: pool}
}

func (r *VoteRepository) Save(ctx context.Context, v domain.Vote) error {
	// Idiomatic SQL: Use named arguments or $ syntax.
	// We cast the int64 timestamp to a Postgres timestamp.
	query := `
        INSERT INTO votes (eviction_id, nominee_id, created_at) 
        VALUES ($1, $2, $3);
    `

	// Convert Unix timestamp to Time object for the driver
	ts := v.Timestamp

	_, err := r.pool.Exec(ctx, query, v.EvictionID, v.NomineeID, ts)
	if err != nil {
		return err
	}

	return nil
}

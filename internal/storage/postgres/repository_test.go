package postgres_test

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	"github.com/jackc/pgx/v5/pgxpool"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/testcontainers/testcontainers-go"
// 	"github.com/testcontainers/testcontainers-go/modules/postgres"
// 	"github.com/testcontainers/testcontainers-go/wait"

// 	"b3e/internal/domain"
// 	storagePg "b3e/internal/storage/postgres"
// )

// func TestVoteRepository_Save(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping integration test")
// 	}

// 	ctx := context.Background()

// 	// 1. Spin up Postgres Container
// 	pgContainer, err := postgres.RunContainer(ctx,
// 		testcontainers.WithImage("postgres:15-alpine"),
// 		postgres.WithDatabase("testdb"),
// 		postgres.WithUsername("user"),
// 		postgres.WithPassword("password"),
// 		testcontainers.WithWaitStrategy(
// 			wait.ForLog("database system is ready to accept connections").
// 				WithOccurrence(2).
// 				WithStartupTimeout(5*time.Second)),
// 	)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer pgContainer.Terminate(ctx)

// 	// 2. Get Connection String
// 	connStr, _ := pgContainer.ConnectionString(ctx, "sslmode=disable")

// 	// 3. Connect & Migrate (Create Table)
// 	pool, err := pgxpool.New(ctx, connStr)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer pool.Close()

// 	_, err = pool.Exec(ctx, `CREATE TABLE votes (id TEXT PRIMARY KEY, candidate TEXT, voter_id TEXT, created_at TIMESTAMP)`)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// 4. Run the Test
// 	repo := storagePg.NewVoteRepository(pool)

// 	vote := domain.Vote{
// 		ID:        "vote_unique_1",
// 		Candidate: "Gopher",
// 		VoterID:   "dev_1",
// 		Timestamp: time.Now().Unix(),
// 	}

// 	// Attempt Save
// 	err = repo.Save(ctx, vote)
// 	assert.NoError(t, err)

// 	// Verify Record Exists
// 	var count int
// 	pool.QueryRow(ctx, "SELECT count(*) FROM votes WHERE id=$1", "vote_unique_1").Scan(&count)
// 	assert.Equal(t, 1, count)
// }

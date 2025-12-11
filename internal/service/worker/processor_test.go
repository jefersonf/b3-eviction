package worker_test

import (
	"b3e/internal/core/command"
	"b3e/internal/domain"
	"b3e/internal/service/worker"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

// MockRepo simulates the database
type MockRepo struct {
	mock.Mock
}

func (r *MockRepo) Save(ctx context.Context, v domain.Vote) error {
	return r.Called(ctx, v).Error(0)
}

func TestWorkerProcess(t *testing.T) {
	// 1. Setup
	mockRepo := new(MockRepo)
	voteChan := make(chan command.CastVote, 1) // buffering to prevent blocking test

	// The command we expect to process
	cmd := command.CastVote{
		NomineeID: "t0",
		Timestamp: time.Now().UTC(),
	}

	// 2. Expectation: the repo SHOULD be called when a message arrives
	mockRepo.On("Save", mock.Anything, mock.MatchedBy(func(v domain.Vote) bool {
		return v.NomineeID == "t0"
	})).Return(nil)

	// 3. Initialize Worker
	w := worker.NewProcessor(mockRepo)

	// 4. Simulating worker logic munually for the test
	ctx, cancel := context.WithCancel(context.Background())

	// Inject the message
	voteChan <- cmd

	// Run worker in background
	go func() {
		w.Run(ctx, voteChan)
	}()

	// Allow some time for processing
	time.Sleep(30 * time.Millisecond)
	cancel()

	// 5. Assert
	mockRepo.AssertExpectations(t)
}

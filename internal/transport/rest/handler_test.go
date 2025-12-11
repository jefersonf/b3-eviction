package rest_test

import (
	"b3e/internal/core/command"
	"b3e/internal/transport/rest"
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBus is our dependency mock.
type MockBus struct {
	mock.Mock
}

func (b *MockBus) Publish(ctx context.Context, cmd command.CastVote) error {
	args := b.Called(ctx, cmd)
	return args.Error(0)
}

func TestIngestVote(t *testing.T) {
	tests := []struct {
		name           string
		payload        string
		mockReturnErr  error
		expectedStatus int
	}{
		{
			name:           "Valid vote returns 202 Accepted",
			payload:        `{"nominee": "abc100"}`,
			mockReturnErr:  nil,
			expectedStatus: http.StatusAccepted,
		},
		{
			name:           "Bus failure returns 500",
			payload:        `{"nominee": "abc100"}`,
			mockReturnErr:  errors.New("queue full"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// 1. Setup mock
			mockBus := new(MockBus)
			if tc.expectedStatus != http.StatusBadRequest {
				mockBus.On("Publish", mock.Anything, mock.AnythingOfType("command.CastVote")).Return(tc.mockReturnErr)
			}
			// 2. Setup handler
			handler := rest.NewVoteHandler(mockBus)
			req := httptest.NewRequest(http.MethodPost, "/votes", bytes.NewBufferString(tc.payload))
			w := httptest.NewRecorder()

			// 3. Execute
			handler.HandleVote(w, req)

			// 4. Assert
			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

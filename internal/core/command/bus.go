package command

import "context"

// Publisher is used by the voting-api.
type Publisher interface {
	Publish(ctx context.Context, cmd CastVote) error
}

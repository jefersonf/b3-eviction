package command

import "context"

// Publisher is used by the Ingestor.
type Publisher interface {
	Publish(ctx context.Context, cmd CastVote) error
}

// Subscriber is used by the Worker.
type Subscriber interface {
	Subscribe(ctx context.Context) (<-chan CastVote, error)
}

// Bus aggregates both Publisher and Subscriber interfaces.
// It decouples the API from the actual queue.
type Bus interface {
	Publisher
	Subscriber
}

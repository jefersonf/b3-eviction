package worker

import (
	"b3e/internal/core/command"
	"b3e/internal/domain"
	"context"
	"log"
)

type VoteRepository interface {
	Save(ctx context.Context, v domain.Vote) error
}

type Processor struct {
	repo VoteRepository
}

func NewProcessor(r VoteRepository) *Processor {
	return &Processor{repo: r}
}

func (p *Processor) Run(ctx context.Context, votes <-chan command.CastVote) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Working shutting down...")
			return
		case cmd, ok := <-votes:
			if !ok {
				return // channel closed
			}
			p.process(ctx, cmd)
		}
	}
}

func (p *Processor) process(ctx context.Context, cmd command.CastVote) {
	// Map Command to Domain entity
	vote := domain.Vote{
		NomineeID: cmd.NomineeID,
		Timestamp: cmd.Timestamp,
	}

	// Attempt to save
	if err := p.repo.Save(ctx, vote); err != nil {
		log.Printf("Failed to save vote %s: %v\n", cmd.NomineeID, err)
	}
}

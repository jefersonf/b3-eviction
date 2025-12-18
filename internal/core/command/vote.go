package command

import (
	"time"
)

// CastVote is the Command to cast a vote.
type CastVote struct {
	EvictionID string
	NomineeID  string
	Timestamp  time.Time
}

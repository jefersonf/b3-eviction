package command

import (
	"time"
)

// CastVote is the command to cast a vote.
type CastVote struct {
	EvictionID string
	NomineeID  string
	Timestamp  time.Time
}

package command

import (
	"time"
)

// CastVote is the Command pattern.
type CastVote struct {
	NomineeID string
	Timestamp time.Time
}

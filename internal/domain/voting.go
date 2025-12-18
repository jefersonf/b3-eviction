package domain

import "time"

type Vote struct {
	EvictionID string
	NomineeID  string
	Timestamp  time.Time
}

type VoteCount struct {
	NomineeID string
	Votes     int64
}

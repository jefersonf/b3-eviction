package domain

import (
	"time"
)

type TimelyStat struct {
	Timedate  time.Time `json:"timedate"`
	NomineeID string    `json:"nominee_id"`
	Votes     int64     `json:"votes"`
}

type EvictionStats struct {
	TotalVotes     int64       `json:"total_votes"`
	VotesByNominee []VoteCount `json:"nominee_votes"`
	Evicted        VoteCount   `json:"evicted"`
}

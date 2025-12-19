package rest

type VoteRequest struct {
	EvictionID string `json:"eviction_id"`
	NomineeID  string `json:"nominee_id"`
}

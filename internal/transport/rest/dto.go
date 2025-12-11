package rest

import "b3e/internal/domain"

type VoteRequest struct {
	NomineeID string `json:"nominee_id"`
}

// ToDomain map DTO to Domain
func (r *VoteRequest) ToDomain() domain.Vote {
	return domain.Vote{NomineeID: r.NomineeID}
}

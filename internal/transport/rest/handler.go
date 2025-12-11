package rest

import (
	"b3e/internal/core/command"
	"encoding/json"
	"net/http"
	"time"
)

type VoteHandler struct {
	bus command.Publisher
}

func NewVoteHandler(b command.Publisher) *VoteHandler {
	return &VoteHandler{bus: b}
}

func (h *VoteHandler) HandleVote(w http.ResponseWriter, r *http.Request) {

	var req VoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	cmd := command.CastVote{NomineeID: req.NomineeID, Timestamp: time.Now().UTC()}

	if err := h.bus.Publish(r.Context(), cmd); err != nil {
		// log error here
		http.Error(w, "Failed to ingest vote", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

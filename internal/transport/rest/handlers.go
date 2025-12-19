package rest

import (
	"b3e/internal/core/command"
	"cmp"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

var DefaultInstanceName = cmp.Or(os.Getenv("INSTANCE_NAME"), "unnamed")

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

	cmd := command.CastVote{
		EvictionID: req.EvictionID,
		NomineeID:  req.NomineeID,
		Timestamp:  time.Now().UTC(),
	}

	if err := h.bus.Publish(r.Context(), cmd); err != nil {
		http.Error(w, "Failed to ingest vote", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusAccepted)
}

func HandleHealthCheck(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"status\": \"up and running\", \"instance\": \"" + DefaultInstanceName + "\"}\n"))
}

func HandleResourceNotFound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("{\"message\": \"resource not found\"}"))
}

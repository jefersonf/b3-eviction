package rest

import (
	"b3e/internal/service/analytics"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type VoteAnalyticsHandler struct {
	stats       analytics.StatsProvider
	timelyStats analytics.TimelyStatsProvider
}

func NewVoteAnalyticsHandler(stats analytics.StatsProvider, timelyStats analytics.TimelyStatsProvider) *VoteAnalyticsHandler {
	return &VoteAnalyticsHandler{
		stats:       stats,
		timelyStats: timelyStats,
	}
}

// HandleHourlyStats ...
func (h *VoteAnalyticsHandler) HandleHourlyStats(w http.ResponseWriter, r *http.Request) {
	timelyStats, err := h.timelyStats.HourlyStats(context.Background())
	if err != nil {
		errMsg := "request for hourly voting statistics could not be processed"
		http.Error(w, errMsg, http.StatusInternalServerError)
		log.Printf("Error: %s\n", err)
		return
	}

	respBody, _ := json.Marshal(timelyStats)
	hash := sha256.Sum256(respBody)
	etag := fmt.Sprintf(`"%x"`, hash)

	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("ETag", etag)

	if match := r.Header.Get("If-None-Match"); match == etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	w.Write(respBody)
}

// HandleMinutelylyStats ...
func (h *VoteAnalyticsHandler) HandleMinutelylyStats(w http.ResponseWriter, r *http.Request) {
	timelyStats, err := h.timelyStats.MinutelyStats(context.Background())
	if err != nil {
		errMsg := "request for minutely voting statistics could not be processed"
		http.Error(w, errMsg, http.StatusInternalServerError)
		log.Printf("Error: %s\n", err)
		return
	}
	json.NewEncoder(w).Encode(timelyStats)
}

// HandleEvictionStats ...
func (h *VoteAnalyticsHandler) HandleEvictionStats(w http.ResponseWriter, r *http.Request) {
	evictionID := r.PathValue("evictionId")
	if len(evictionID) == 0 {
		errMsg := "missing required eviction id"
		http.Error(w, errMsg, http.StatusBadRequest)
		log.Printf("Error: %s\n", errMsg)
		return
	}
	evictionStats, err := h.stats.EvictionStats(r.Context(), evictionID)
	if err != nil {
		errMsg := fmt.Sprintf("request for eviction id:%s statistics could not be processed", evictionID)
		http.Error(w, errMsg, http.StatusInternalServerError)
		log.Printf("Error: %s\n", errMsg)
	}

	json.NewEncoder(w).Encode(evictionStats)
}

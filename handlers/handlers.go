package handlers

import (
	"encoding/json"
	"github.com/skylibdrvlz/20.11.2025-links-checker/checker"
	"github.com/skylibdrvlz/20.11.2025-links-checker/models"
	"github.com/skylibdrvlz/20.11.2025-links-checker/storage"
	"log/slog"
	"net/http"
)

type Handler struct {
	storage *storage.Storage
	checker *checker.Checker
}

func NewHandler(storage *storage.Storage) *Handler {
	return &Handler{
		storage: storage,
		checker: checker.NewChecker(),
	}
}

func (h *Handler) CheckLinks(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	var req models.CheckLinksRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	results := h.checker.CheckLinks(r.Context(), req.Links)
	linksNum := h.storage.SaveLinkSet(results)

	resp := models.CheckLinksResponse{
		Links:    results,
		LinksNum: linksNum,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

	slog.Info("Links checked", "count", len(req.Links), "id", linksNum)
}

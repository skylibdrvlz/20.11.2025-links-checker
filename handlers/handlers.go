package handlers

import (
	"encoding/json"
	"github.com/skylibdrvlz/20.11.2025-links-checker/checker"
	"github.com/skylibdrvlz/20.11.2025-links-checker/models"
	"github.com/skylibdrvlz/20.11.2025-links-checker/pdf"
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

func (h *Handler) GenerateReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	var req models.LinksListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Invalid JSON in generate report request", "error", err)
		http.Error(w, "Invalid JSON", 400)
		return
	}

	slog.Info("Generating PDF report", "requested_ids", req.LinksList)

	linkSets, err := h.storage.GetLinkSets(req.LinksList)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	if len(linkSets) == 0 {
		http.Error(w, "No link sets found", 404)
		return
	}

	pdfData, err := pdf.GeneratePDF(linkSets)
	if err != nil {
		slog.Error("PDF generation failed", "error", err, "link_sets_count", len(linkSets))
		http.Error(w, "Error generating PDF", 500)
		return
	}

	slog.Info("PDF report generated successfully",
		"link_sets_count", len(linkSets),
		"pdf_size_bytes", len(pdfData))

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=report.pdf")
	w.Write(pdfData)
}

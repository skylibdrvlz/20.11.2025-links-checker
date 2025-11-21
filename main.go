package main

import (
	"github.com/skylibdrvlz/20.11.2025-links-checker/handlers"
	"github.com/skylibdrvlz/20.11.2025-links-checker/storage"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	storage := storage.NewStorage("data.json")
	handler := handlers.NewHandler(storage)

	mux := http.NewServeMux()
	mux.HandleFunc("/check-links", handler.CheckLinks)
	mux.HandleFunc("/generate-report", handler.GenerateReport)

	slog.Info("Server starting", "port", "8080")
	http.ListenAndServe(":8080", mux)
}

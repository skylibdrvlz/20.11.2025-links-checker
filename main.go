package main

import (
	"github.com/skylibdrvlz/20.11.2025-links-checker/handlers"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	handler := handlers.NewHandler()

	mux := http.NewServeMux()
	mux.HandleFunc("/check-links", handler.CheckLinks)

	slog.Info("Server starting", "port", "8080")
	http.ListenAndServe(":8080", mux)
}

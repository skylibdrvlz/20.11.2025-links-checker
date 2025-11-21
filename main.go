package main

import (
	"context"
	"github.com/skylibdrvlz/20.11.2025-links-checker/handlers"
	"github.com/skylibdrvlz/20.11.2025-links-checker/storage"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	storage := storage.NewStorage("data.json")
	handler := handlers.NewHandler(storage)

	mux := http.NewServeMux()
	mux.HandleFunc("/check-links", handler.CheckLinks)
	mux.HandleFunc("/generate-report", handler.GenerateReport)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		slog.Info("Graceful shutdown activated")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			slog.Error("Graceful shutdown failed", "error", err)
		}
	}()

	slog.Info("Server starting", "port", "8080")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("Server failed", "error", err)
	}

	slog.Info("Server stopped")
}

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
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)

	go func() {
		slog.Info("server started", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	<-quit
	slog.Info("graceful shutdown initiated")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("shutdown error", "err", err)
	} else {
		slog.Info("server gracefully stopped")
	}

}

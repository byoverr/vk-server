package main

import (
	"backend/internal/config"
	"backend/internal/handlers"
	"backend/internal/storage"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	ctx := context.Background()

	urlAddress := fmt.Sprintf("%s:%d", cfg.HTTPServer.Host, cfg.HTTPServer.Port)

	repo, err := storage.PostgresqlOpen(cfg, ctx)
	if err != nil {
		log.Fatal(err)
	}

	handler := handlers.NewHandler(repo, logger)
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)
	handlers.InitRoutes(router, handler)

	logger.Info("start server", slog.String("address", urlAddress))
	srv := &http.Server{
		Addr:    urlAddress,
		Handler: router,
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Error("failed to start")
	}
}

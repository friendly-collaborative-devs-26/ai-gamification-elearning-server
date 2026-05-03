package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ai-gamification-elearning-server/pkg/config"
	"ai-gamification-elearning-server/pkg/logger"

	"github.com/gin-gonic/gin"
)

var buildVersion = "dev"

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("critical: failed to load config: %v", err)
	}

	if err := logger.Init(cfg.Logger); err != nil {
		log.Fatalf("critical: failed to initialise logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("ai-gamification-elearning-server starting",
		logger.String("env", cfg.App.Env),
		logger.String("version", buildVersion),
	)

	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := setupRouter(cfg)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.App.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("http server listening",
			logger.String("addr", srv.Addr),
			logger.String("url", fmt.Sprintf("http://localhost:%d/api/%s", cfg.App.Port, cfg.App.Version)),
		)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("http server error", logger.Err(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutdown signal received — draining in-flight requests...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("forced shutdown — context deadline exceeded", logger.Err(err))
	}

	logger.Info("server stopped cleanly")
}

func setupRouter(cfg *config.Config) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"time":    time.Now().UTC().Format(time.RFC3339),
			"version": buildVersion,
		})
	})

	apiVersion := fmt.Sprintf("/api/%s", cfg.App.Version)
	api := router.Group(apiVersion)
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		})
	}

	return router
}

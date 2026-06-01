package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"ai-gamification-elearning-server/internal/infrastructure/persistence"
	"ai-gamification-elearning-server/pkg/config"
	"ai-gamification-elearning-server/pkg/logger"
)

func main() {
	if err := chdirToProjectRoot(); err != nil {
		log.Fatalf("migrate: failed to set working directory: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("migrate: failed to load config: %v", err)
	}

	if err := logger.Init(cfg.Logger); err != nil {
		log.Fatalf("migrate: failed to initialise logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("running database migrations",
		logger.String("host", cfg.Database.Host),
		logger.String("user", cfg.Database.User),
		logger.String("db", cfg.Database.Name),
		logger.String("env", cfg.App.Env),
	)

	if _, err := persistence.NewDB(cfg, logger.Get()); err != nil {
		logger.Fatal("migrate: failed", logger.Err(err))
	}

	logger.Info("migrations completed successfully")
}

func chdirToProjectRoot() error {
	if root := os.Getenv("PROJECT_ROOT"); root != "" {
		return os.Chdir(root)
	}

	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return os.Chdir(dir)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return fmt.Errorf("go.mod not found — set PROJECT_ROOT env var manually")
		}
		dir = parent
	}
}

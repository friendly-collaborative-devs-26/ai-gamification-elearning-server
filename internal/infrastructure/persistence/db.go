package persistence

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"ai-gamification-elearning-server/internal/infrastructure/persistence/models"
	"ai-gamification-elearning-server/pkg/config"
)

func NewDB(cfg *config.Config, log *zap.Logger) (*gorm.DB, error) {
	dsn := buildDSN(cfg)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(resolveLogLevel(cfg)),
	})
	if err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}

	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("running migrations: %w", err)
	}

	if err := configurePool(db, cfg); err != nil {
		return nil, fmt.Errorf("configuring connection pool: %w", err)
	}

	log.Info("database connected", zap.String("host", cfg.Database.Host))
	return db, nil
}

func buildDSN(cfg *config.Config) string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s&TimeZone=UTC",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)
}

func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.UserModel{},
		&models.PermissionModel{},
		&models.StoreItemModel{},
		&models.AccountModel{},
		&models.UserPermissionModel{},
		&models.BlockModel{},
		&models.ExerciseModel{},
		&models.UserScoreModel{},
		&models.ExerciseStatusModel{},
		&models.ExerciseExecutionModel{},
		&models.BuyedStoreItemModel{},
	)
}

func configurePool(db *gorm.DB, cfg *config.Config) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetimeSecs) * time.Second)
	return nil
}

func resolveLogLevel(cfg *config.Config) gormlogger.LogLevel {
	if cfg.App.Env == "development" {
		return gormlogger.Info
	}
	return gormlogger.Warn
}

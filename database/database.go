package database

import (
	"fmt"
	"log"
	"os"
	"project/config"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB(cfg config.Config) (*gorm.DB, error) {
	// Configure the database logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Enable color output
		},
	)

	db, err := gorm.Open(postgres.Open(makePostgresString(cfg)), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	db.Exec(`
		DO $$ BEGIN CREATE TYPE orderstatus AS ENUM('created', 'processed', 'canceled', 'completed');
		EXCEPTION WHEN duplicate_object THEN null; END $$;
	`)

	// Call Migrate function to auto-migrate database schemas
	if cfg.DBMigrate {
		err = Migrate(db)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	// Call See function to auto-migrate database schemas
	if cfg.DBSeeding {
		err = SeedAll(db)
	}
	if err != nil {
		return nil, err
	}

	return db, nil
}

// makePostgresString creates the PostgreSQL DSN (Data Source Name)
func makePostgresString(cfg config.Config) string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName, cfg.DBPassword)
}

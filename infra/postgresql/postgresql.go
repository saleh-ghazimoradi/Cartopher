package postgresql

import (
	"fmt"

	"github.com/saleh-ghazimoradi/Cartopher/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresql(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC", cfg.Postgresql.Host, cfg.Postgresql.User, cfg.Postgresql.Password, cfg.Postgresql.Name, cfg.Postgresql.Port, cfg.Postgresql.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgresql database: %w", err)
	}
	return db, nil
}

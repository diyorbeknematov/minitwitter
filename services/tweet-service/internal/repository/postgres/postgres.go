package postgres

import (
	"fmt"

	"github.com/diyorbeknematov/minitwitter/services/tweet-service/internal/config"
	"github.com/jmoiron/sqlx"
)

func ConnectionDB(cfg *config.Config) (*sqlx.DB, error) {
	conn := fmt.Sprintf("host=%s  port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := sqlx.Connect("postgres", conn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

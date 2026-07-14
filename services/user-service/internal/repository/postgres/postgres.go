package postgres

import (
	"fmt"

	"github.com/diyorbek/minitwitter/services/user-service/internal/config"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

func DBConnection(cfg *config.Config) (*sqlx.DB, error) {
	conn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sqlx.Connect("postgres", conn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

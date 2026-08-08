package postgres

import (
	"fmt"

	"github.com/diyorbeknematov/minitwitter/services/tweet-service/internal/config"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

func ConnectionDB(cfg config.DBConfig) (*sqlx.DB, error) {
	conn := fmt.Sprintf("host=%s  port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name,
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

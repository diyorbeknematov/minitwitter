package postgres

import (
	"fmt"

	"github.com/diyorbek/minitwitter/services/user-service/internal/config"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

func DBConnection(cfg config.DBConfig) (*sqlx.DB, error) {
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)

	db, err := sqlx.Connect("postgres", conn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

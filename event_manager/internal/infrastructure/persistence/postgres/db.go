package postgres

import (
	"database/sql"
	"fmt"

	"github.com/XRS0/ToTalkB/event_manager/config"
	_ "github.com/lib/pq"
)

func NewPostgresDB(cfg config.Database) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	return sql.Open("postgres", connStr)
}

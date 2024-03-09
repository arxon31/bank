package postgres

import (
	"database/sql"
	"log/slog"

	_ "github.com/jackc/pgx/stdlib" // pgx driver
)

type Postgres struct {
	DB     *sql.DB
	Logger *slog.Logger
}

func New(url string, logger *slog.Logger) (*Postgres, error) {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Postgres{DB: db, Logger: logger}, nil
}

func (p *Postgres) Close() error {
	return p.DB.Close()
}

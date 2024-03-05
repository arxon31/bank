package postgres

import (
	"database/sql"
	_ "github.com/jackc/pgx/stdlib"
	"log/slog"
)

type Postgres struct {
	DB     *sql.DB
	logger *slog.Logger
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

	return &Postgres{DB: db, logger: logger}, nil
}

func (p *Postgres) Close() error {
	return p.DB.Close()
}

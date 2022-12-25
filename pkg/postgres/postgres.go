package postgres

import (
	"database/sql"
	"log"
	"time"

	"golang.org/x/exp/slog"

	_ "github.com/lib/pq"
)

const (
	_defaultConnAttempts = 3
	_defaultConnTimeout  = time.Second
)

type Postgres struct {
	connAttempts int
	connTimeout  time.Duration

	DB *sql.DB
}

func NewPostgresDB(url string, opts ...Option) (*Postgres, error) {
	slog.Info("CONN", "connect string", url)

	pg := &Postgres{
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	for _, opt := range opts {
		opt(pg)
	}

	var err error
	for pg.connAttempts > 0 {
		pg.DB, err = sql.Open("postgres", url)
		if err != nil {
			break
		}

		log.Printf("Postgres is trying to connect, attempts left: %d", pg.connAttempts)

		time.Sleep(pg.connTimeout)

		pg.connAttempts--
	}

	return pg, nil
}

func (p *Postgres) Close() {
	if p.DB != nil {
		p.DB.Close()
	}
}

package postgres

import "time"

type Option func(*Postgres)

func MaxPoolSize(size int) Option {
	return func(p *Postgres) {
		p.maxPoolSize = size
	}
}

func ConnAttempts(attempts int) Option {
	return func(p *Postgres) {
		p.connAttempts = attempts
	}
}

func ConnTimeout(timeout time.Duration) Option {
	return func(p *Postgres) {
		p.connTimeout = timeout
	}
}

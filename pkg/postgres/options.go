package postgres

import "time"

type Option func(*postgres)

func ConnAttempts(attempts int) Option {
	return func(p *postgres) {
		p.connAttempts = attempts
	}
}

func ConnTimeout(timeout time.Duration) Option {
	return func(p *postgres) {
		p.connTimeout = timeout
	}
}

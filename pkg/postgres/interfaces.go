package postgres

import "database/sql"

type DBEngine interface {
	GetDB() *sql.DB
	Configure(...Option) DBEngine
	Close()
}

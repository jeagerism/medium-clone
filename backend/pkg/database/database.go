package database

import "github.com/jmoiron/sqlx"

type Database interface {
	GetDb() *sqlx.DB
}

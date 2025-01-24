package database

import (
	"fmt"
	"sync"

	"github.com/jeagerism/medium-clone/backend/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Import the postgres driver
)

type postgresDatabase struct {
	DB *sqlx.DB
}

var (
	once       sync.Once
	dbInstance *postgresDatabase
	initErr    error
)

// NewPostgresDatabase initializes a singleton instance of the PostgreSQL database.
func NewPostgresDatabase(conf *config.Config) (*postgresDatabase, error) {
	if conf == nil {
		return nil, fmt.Errorf("database configuration is missing")
	}

	once.Do(func() {
		dsn := conf.Db().Url()
		db, err := sqlx.Open("postgres", dsn)
		if err != nil {
			initErr = fmt.Errorf("failed to open database: %w", err)
			return
		}

		if err = db.Ping(); err != nil {
			initErr = fmt.Errorf("failed to ping database: %w", err)
			return
		}

		dbInstance = &postgresDatabase{DB: db}
	})

	if initErr != nil {
		return nil, initErr
	}

	return dbInstance, nil
}

// GetDb returns the underlying *sqlx.DB instance.
func (p *postgresDatabase) GetDb() *sqlx.DB {
	return p.DB
}

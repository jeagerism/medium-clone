package database

import (
	"fmt"
	"log"
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
)

func NewPostgresDatabase(conf *config.Config) Database {
	once.Do(func() {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
			conf.DB.Host,
			conf.DB.User,
			conf.DB.Password,
			conf.DB.DBName,
			conf.DB.Port,
			conf.DB.SSLMode,
			conf.DB.Timezone,
		)

		db, err := sqlx.Open("postgres", dsn)
		if err != nil {
			log.Fatalf("failed to open database: %v", err)
		}

		if err = db.Ping(); err != nil {
			log.Fatalf("failed to ping database: %v", err)
		}

		dbInstance = &postgresDatabase{DB: db}
	})

	return dbInstance
}

func (p *postgresDatabase) GetDb() *sqlx.DB {
	return dbInstance.DB
}

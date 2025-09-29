package services

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

// https://donchev.is/post/working-with-postgresql-in-go-using-pgx/

type Postgres struct {
	Db *pgxpool.Pool
}

var (
	pgInstance *Postgres
	pgOnce     sync.Once
)

func NewPG(cfg *DatabaseConfig) (_ *Postgres, err error) {
	pgOnce.Do(func() {
		dsn := toDBS(cfg)
		db, err := pgxpool.New(context.Background(), dsn)

		if err != nil {
			log.Fatal(fmt.Errorf("unable to create connection pool: %v", err))
		}
		pgInstance = &Postgres{db}
	})

	return pgInstance, nil
}

func toDBS(cfg *DatabaseConfig) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)
}

func (pg *Postgres) Close() {
	pg.Db.Close()
}

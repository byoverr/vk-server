package storage

import (
	"backend/internal/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type PostgresqlDB struct {
	db  *pgxpool.Pool
	ctx context.Context
}

func PostgresqlOpen(cfg *config.Config, ctx context.Context) (*PostgresqlDB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Storage.User, cfg.Storage.Password, cfg.Storage.Host,
		cfg.Storage.Port, cfg.Storage.DbName, cfg.Storage.SSLMode)
	fmt.Println(connStr)
	db, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("Unable to connection to database: %v\n", err)
		return nil, err
	}
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	postgresql := &PostgresqlDB{db: db}

	err = postgresql.Init()

	if err != nil {
		return nil, err
	}

	return postgresql, nil
}

func (s *PostgresqlDB) Init() error {
	q := `
CREATE TABLE IF NOT EXISTS statuses (
    id SERIAL PRIMARY KEY,
    ip TEXT,
    ping_time TEXT NULL,
    last_check timestamp NULL 
);
`

	if _, err := s.db.Exec(context.Background(), q); err != nil {
		return err
	}
	return nil
}

package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func InitDBCConnection(ctx context.Context, pgConfig *pgxpool.Config) (*pgxpool.Pool, error) {
	for i := 0; i < 10; i++ {
		db, err := pgxpool.ConnectConfig(ctx, pgConfig)
		if err != nil {
			log.Printf("connect to db: %s retry to connect", err.Error())
			time.Sleep(5 * time.Second)

			continue
		}

		return db, nil
	}

	return nil, fmt.Errorf("database is unreachable")
}

func ConfigureDBConnection(ctx context.Context) *pgxpool.Config {
	pgConf, err := pgxpool.ParseConfig("postgres://postgres:password@172.17.0.1:8081/midasinvestment")
	if err != nil {
		panic(err)
	}

	pgConf.MaxConnIdleTime = 10 * time.Second
	pgConf.MaxConnLifetime = 10 * time.Second
	pgConf.MaxConns = 100

	return pgConf
}

package database

import (
	"context"
	"embed"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func ConnectPG() *pgxpool.Pool {
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("GOOSE_DBSTRING"))
	if err != nil {
		panic(fmt.Errorf("unable to connect to database: %v", err))
	}

	db := stdlib.OpenDBFromPool(dbpool)

	// Run migrations
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("pgx"); err != nil {
		panic(fmt.Errorf("failed to set dialect: %v", err))
	}

	if err := goose.Up(db, "migrations"); err != nil {
		panic(fmt.Errorf("failed to run migrations: %v", err))
	}

	return dbpool

}

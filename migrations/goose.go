package migrations

import (
	"embed"
	"fmt"
	"go.uber.org/zap"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var embedMigrations embed.FS

type ZapLogger struct {
	*zap.Logger
}

func (z ZapLogger) Fatalf(format string, v ...interface{}) {
	z.Fatal(fmt.Sprintf(format, v...))
}

func (z ZapLogger) Printf(format string, v ...interface{}) {
	z.Info(fmt.Sprintf(format, v...))
}

func RunMigrations(l *zap.Logger, pool *pgxpool.Pool) error {
	goose.SetBaseFS(embedMigrations)
	goose.SetLogger(ZapLogger{l})

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set dialect: %w", err)
	}

	db := stdlib.OpenDBFromPool(pool)
	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}

package main

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	"github.com/huandu/go-sqlbuilder"
	"github.com/illiafox/spy-cat-test-assignment/app/config"
	"github.com/illiafox/spy-cat-test-assignment/app/internal/repository/catapi"
	"github.com/illiafox/spy-cat-test-assignment/app/internal/repository/postgres"
	"github.com/illiafox/spy-cat-test-assignment/app/internal/service"
	"github.com/illiafox/spy-cat-test-assignment/app/internal/transport/http"
	"github.com/illiafox/spy-cat-test-assignment/app/pkg/poolwrapper"
	"github.com/illiafox/spy-cat-test-assignment/migrations"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"os"
	"os/signal"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("read config: %v", err)
	}

	var loggerConfig zap.Config
	if cfg.Debug {
		loggerConfig = zap.NewDevelopmentConfig()
	} else {
		loggerConfig = zap.NewProductionConfig()
	}
	loggerConfig.DisableStacktrace = cfg.DisableStacktrace

	logger, err := loggerConfig.Build()
	if err != nil {
		log.Fatalf("failed to start zap logger: %w", err)
	}

	defer logger.Sync()

	//

	pool, err := pgxpool.New(context.TODO(), cfg.PostgresURI)
	if err != nil {
		logger.Fatal("failed to init pgxpool", zap.Error(err))
	}
	defer pool.Close()

	if err = migrations.RunMigrations(logger, pool); err != nil {
		logger.Fatal("failed to run migrations", zap.Error(err))
	}

	sqlbuilder.DefaultFlavor = sqlbuilder.PostgreSQL
	db := poolwrapper.NewPool(pool)
	catsRepository := postgres.NewCatsRepository(db)
	missionsRepository := postgres.NewMissionsRepository(db)
	targetsRepository := postgres.NewTargetsRepository(db)
	notesRepository := postgres.NewNotesRepository(db)

	catBreedChecker := catapi.NewClient()

	//

	service := service.NewService(
		catBreedChecker,
		catsRepository,
		missionsRepository,
		targetsRepository,
		notesRepository,
		catsRepository,
	)

	//

	server := http.NewServer(cfg, logger, service)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := server.Start(); err != nil {
			logger.Error("failed to start http server", zap.Error(err))
			cancel()
		}
	}()

	ctx, cancel = signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	<-ctx.Done()

	logger.Info("Shutting down server")
}

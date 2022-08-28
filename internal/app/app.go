package app

import (
	"context"
	"os"
	"os/signal"
	"time"

	"gitlab.com/g6834/team31/analytics/internal/adapters/database/postgres"
	"gitlab.com/g6834/team31/analytics/internal/adapters/grpc"
	"gitlab.com/g6834/team31/analytics/internal/adapters/http"
	"gitlab.com/g6834/team31/analytics/internal/adapters/mq"
	"gitlab.com/g6834/team31/analytics/internal/config"
	"gitlab.com/g6834/team31/analytics/internal/domain/analytics"
	"gitlab.com/g6834/team31/auth/pkg/logging"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose"
)

const (
	dbTimeOut = time.Second * 10
)

// Run creates objects via constructors.
func Run(ctx context.Context, cfg *config.Config) {
	ctx, cancel := context.WithCancel(ctx)
	l := logging.New(cfg.Log.Level)
	l.Info().Msgf("config: %+v", cfg)
	pgConnStr := os.Getenv("PG_CONNSTR")
	if pgConnStr != "" {
		l.Info().Msg("found address to postgres in env")
		cfg.PG.URL = pgConnStr
	}

	l.Info().Msg("connecting to db")
	dbCtx, _ := context.WithTimeout(ctx, dbTimeOut)
	db, err := postgres.New(dbCtx, cfg.PG)
	if err != nil {
		l.Fatal().Msg(err.Error())
	}

	// Migrations
	l.Info().Msg("starting migrations")
	pgxCfg, err := pgx.ParseConfig(cfg.PG.URL)
	if err != nil {
		l.Fatal().Msgf("pgx.ParseConfig: %+v", err)
	}
	pgxCfg.PreferSimpleProtocol = true

	gooseDb := stdlib.OpenDB(*pgxCfg)
	defer func() {
		db.Close()
	}()
	err = goose.SetDialect("postgres")
	if err != nil {
		l.Fatal().Msg(err.Error())
	}
	if err := goose.Up(gooseDb, "migrations"); err != nil {
		l.Fatal().Msg(err.Error())
	}

	// grpc Auth
	l.Info().Msg("starting grpc auth client")
	authClient, err := grpc.New(ctx, cfg.Auth.HOST, cfg.Auth.PORT, &l)
	if err != nil {
		l.Fatal().Msg(err.Error())
	}

	// analytic service
	analyticService := analytics.New(db, &l)
	server := http.New(analyticService, authClient, cfg, &l)

	l.Info().Msg("kafka client starting")
	mq := mq.BuildMQClient(cfg.Kafka, analyticService, &l)
	errChan := mq.RunMQ(ctx)

	l.Info().Msg("Task grpc server starting")
	tasksGRPCChanErr := grpc.LaunchGRPCServer(cfg.TasksEventsQueue.PORT, analyticService, &l)

	l.Info().Msg("HTTP server starting")
	httpErrChan := server.Start(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	select {
	case err := <-errChan:
		l.Error().Err(err).Msg("kafka problem")
	case err := <-httpErrChan:
		l.Error().Err(err).Msg("http problem")
	case err := <-tasksGRPCChanErr:
		l.Error().Err(err).Msg("grpc tasks problem")
	case <-c:
		l.Info().Msg("shutting down the app")
	}
	cancel() // точно ли все успевает закрыться?
}

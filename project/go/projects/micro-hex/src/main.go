package main

import (
	"context"
	"hesp_channel/ports/config"
	"hesp_channel/ports/repository/mssql"
	"hesp_channel/ports/worker"
	"hesp_channel/service"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-redis/redis/v9"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	zap_logger "hesp_channel/ports/logger"
)

var (
	cfg    *service.Config
	logger service.LogService
)

func init() {
	var err error
	cfg, err = config.NewConfigService().LoadConfig()
	if err != nil {
		panic(err)
	}

	logger, err = getLogService()
	if err != nil {
		panic(err)
	}
}

func main() {
	var (
		err error
	)

	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	services, err := initializeService(mainCtx)
	if err != nil {
		logger.Errorf("context=\"setup\" error=\"%s\"", err)
		os.Exit(1)
	}

	g, gCtx := errgroup.WithContext(mainCtx)

	for _, s := range services {
		backgroundService := s
		g.Go(func() error {
			return backgroundService.Start(gCtx)
		})
	}

	logger.Info("status=Started")

	err = g.Wait()
	if err != nil {
		logger.Errorf("process=\"Error group\" error=\"%s\"", err)
		os.Exit(2)
	}
	logger.Info("message=\"Service shutting down\"")

}

func initializeService(ctx context.Context) ([]service.BackgroundService, error) {

	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    cfg.Redis.Addresses,
		Password: cfg.Redis.Password,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, errors.Wrap(err, "ping redis")
	}

	sqlCm := mssql.NewSqlConnectionManager(cfg.Sql)
	hespDataRepo := mssql.NewHespDataRepository(sqlCm)

	chs := service.NewHespChannelService(cfg.Hesp, hespDataRepo, logger)

	return []service.BackgroundService{
		worker.NewStreamEventConsumer(cfg, rdb, logger, chs),
	}, nil

}

func getLogService() (service.LogService, error) {
	zapConfig, err := zap_logger.BuildConfig(cfg)
	if err != nil {
		return nil, err
	}
	return zap_logger.NewZapLogger(zapConfig)
}

package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/iftsoft/linker/client/callback"
	"github.com/iftsoft/linker/handler"

	"github.com/iftsoft/driver/config"
	"github.com/iftsoft/driver/system"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	slog.SetDefault(logger)
	logger.Info("Start equipment application")

	if err := Run(ctx, logger); err != nil {
		logger.Error("Equipment application error", "error", err)
	}

	logger.Info("Equipment application is stopped")
}

func Run(ctx context.Context, log *slog.Logger) error {
	// callback client init
	callbackCli, err := callback.NewCallbackClient(ctx, log, "127.0.0.1:9090")
	if err != nil {
		return fmt.Errorf("callback client failed: %w", err)
	}
	defer callbackCli.Close()

	appCfg := &config.AppConfig{}
	sysDriver := system.NewSystemDriver(log, appCfg, callbackCli)

	// gRPC Server init
	address := "127.0.0.1:9098"
	grpcSrv := handler.NewManagerServer(log, address, sysDriver)
	if grpcSrv == nil {
		return errors.New("device server is nil")
	}
	defer grpcSrv.Shutdown()

	// gRPC Server start
	go func() {
		if err := grpcSrv.Start(); err != nil {
			log.Error("GRPC server start failed", "error", err)
		}
	}()

	log.Info("Application is running now, press Ctrl+C to shutdown")
	<-ctx.Done()

	return nil
}

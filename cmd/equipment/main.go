package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/iftsoft/driver/config"
	"github.com/iftsoft/driver/device"
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

	appCfg := &config.AppConfig{}
	creator := &device.DummyCreator{}
	app := system.NewApplication(logger,
		system.WithAppConfig(appCfg),
		system.WithAppDeviceCreator(creator),
	)
	if err := app.Run(ctx); err != nil {
		logger.Error("Equipment application error", "error", err)
	}

	logger.Info("Equipment application is stopped")
}

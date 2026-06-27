package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/iftsoft/driver/config"
	"github.com/iftsoft/driver/device"
	"github.com/iftsoft/driver/system"
)

type AppParams struct {
	AppName string
	DevName string
	LogPath string
	CfgPath string
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	appCfg := &config.AppConfig{}
	logger, err := system.RunBootstrap(appCfg)
	if err != nil {
		fmt.Printf("Failed to run app bootstrap: %s\n", err)
		return
	}
	logger.Info("Start equipment application")

	creator := &device.DummyCreator{}
	app := system.NewApplication(logger,
		system.WithAppConfig(appCfg),
		system.WithAppDeviceCreator(creator),
	)
	if err = app.Run(ctx); err != nil {
		logger.Error("Equipment application error", "error", err)
	}

	logger.Info("Equipment application is stopped")
}

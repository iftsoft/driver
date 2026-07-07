package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

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

	setup, err := system.RunBootstrap()
	if err != nil {
		fmt.Printf("Failed to run app bootstrap: %s\n", err)
		return
	}
	log := setup.Logger
	log.Info("Start equipment application")

	creator := &device.DummyCreator{}
	app := system.NewApplication(setup, creator)
	if err = app.Run(ctx); err != nil {
		log.Error("Equipment application error", "error", err)
	}

	log.Info("Equipment application is stopped")
}

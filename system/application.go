package system

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/iftsoft/linker/client/callback"
	"github.com/iftsoft/linker/handler"

	"github.com/iftsoft/driver/config"
	"github.com/iftsoft/driver/device"
)

type Application struct {
	log     *slog.Logger
	config  *config.AppConfig
	creator device.DeviceCreator
}

type AppOption func(*Application)

func NewApplication(log *slog.Logger, opts ...AppOption) *Application {
	app := &Application{log: log}
	for _, opt := range opts {
		opt(app)
	}
	return app
}

func WithAppConfig(config *config.AppConfig) AppOption {
	return func(app *Application) {
		app.config = config
	}
}

func WithAppDeviceCreator(creator device.DeviceCreator) AppOption {
	return func(app *Application) {
		app.creator = creator
	}
}

func (app *Application) Run(ctx context.Context) error {
	// callback client init
	callbackCli, err := callback.NewCallbackClient(ctx, app.log, "127.0.0.1:9090")
	if err != nil {
		return fmt.Errorf("callback client failed: %w", err)
	}
	defer callbackCli.Close()

	devDriver := NewDeviceDriver(app.log,
		WithDeviceConfig(&app.config.Device),
		WithDeviceCallback(callbackCli),
		WithDeviceCreator(app.creator),
	)
	sysDriver := NewSystemDriver(app.log,
		WithSystemCallback(callbackCli),
		WithDeviceDriver(devDriver),
	)

	// gRPC Server init
	address := "127.0.0.1:9098"
	grpcSrv := handler.NewManagerServer(app.log, address, sysDriver)
	if grpcSrv == nil {
		return errors.New("device server is nil")
	}
	defer grpcSrv.Shutdown()

	// gRPC Server start
	go func() {
		if err = grpcSrv.Start(); err != nil {
			app.log.Error("GRPC server start failed", "error", err)
		}
	}()

	app.log.Info("Application is running now, press Ctrl+C to shutdown")
	<-ctx.Done()

	return nil
}

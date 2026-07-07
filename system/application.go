package system

import (
	"context"
	"errors"
	"fmt"

	"github.com/iftsoft/linker/client/callback"
	"github.com/iftsoft/linker/handler"

	"github.com/iftsoft/driver/device"
)

type Application struct {
	setup   *AppSetup
	creator device.DeviceCreator
}

func NewApplication(setup *AppSetup, creator device.DeviceCreator) *Application {
	app := &Application{
		setup:   setup,
		creator: creator,
	}
	return app
}

func (app *Application) Run(ctx context.Context) error {
	// callback client init
	callbackCli, err := callback.NewCallbackClient(ctx, app.setup.Logger, "127.0.0.1:9090")
	if err != nil {
		return fmt.Errorf("callback client failed: %w", err)
	}
	defer callbackCli.Close()

	devDriver := NewDeviceDriver(app.setup, callbackCli, app.creator)
	sysDriver := NewSystemDriver(app.setup.Logger, callbackCli, devDriver)

	// gRPC Server init
	address := "127.0.0.1:9098"
	grpcSrv := handler.NewManagerServer(app.setup.Logger, address, sysDriver)
	if grpcSrv == nil {
		return errors.New("device server is nil")
	}
	defer grpcSrv.Shutdown()

	// gRPC Server start
	go func() {
		if err = grpcSrv.Start(); err != nil {
			app.setup.Logger.Error("GRPC server start failed", "error", err)
		}
	}()

	app.setup.Logger.Info("Application is running now, press Ctrl+C to shutdown")
	<-ctx.Done()

	return nil
}

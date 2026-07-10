package system

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/iftsoft/linker/client/callback"
	"github.com/iftsoft/linker/handler"
	"github.com/iftsoft/linker/model"

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
	port, err := getFreePort()
	if err != nil {
		return fmt.Errorf("get free port failed: %w", err)
	}

	address := fmt.Sprintf("127.0.0.1:%d", port)
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

	params := app.setup.Params
	greeting := &model.GreetingInfo{
		AppName:  params.AppName,
		DevName:  params.DevName,
		GrpcPort: int64(port),
	}
	err = callbackCli.GreetingInfo(ctx, greeting)

	app.setup.Logger.Info("Application is running now, press Ctrl+C to shutdown")
	<-ctx.Done()

	return nil
}

func getFreePort() (int, error) {
	// Bind to port 0 on localhost to let the OS choose a free port
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = l.Close() // Free the port right after finding it
	}()

	// Cast the address to access the specific Port field
	return l.Addr().(*net.TCPAddr).Port, nil
}

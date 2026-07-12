package device

import (
	"context"
	"log/slog"

	"github.com/iftsoft/linker/model"

	"github.com/iftsoft/driver/config"
)

type Callback interface {
	model.DeviceCallback
	model.PrinterCallback
	model.ReaderCallback
	model.ValidatorCallback
}

type DeviceWorker interface {
	DeviceSettings() model.DeviceSetup
	//InitDevice(ctx context.Context) error
	StartDevice(ctx context.Context, query *model.ConfigUpdate) error
	StopDevice(ctx context.Context) error
	CheckDevice(ctx context.Context) (*model.DeviceMetrics, error)
	DeviceTimer(ctx context.Context, unix int64) error
}

type CreatorParams struct {
	DevName  string
	Logger   *slog.Logger
	Config   *config.DeviceConfig
	Callback Callback
}

type DeviceCreator interface {
	CreateDevice(params CreatorParams) (any, error)
}

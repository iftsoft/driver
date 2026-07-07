package device

import (
	"context"
	"log/slog"

	"github.com/iftsoft/linker/model"

	"github.com/iftsoft/driver/config"
)

type Callback interface {
	//model.SystemCallback
	model.DeviceCallback
	model.PrinterCallback
	model.ReaderCallback
	model.ValidatorCallback
}

type DeviceWorker interface {
	InitDevice(ctx context.Context) error
	StartDevice(ctx context.Context, query *model.SystemConfig) error
	StopDevice(ctx context.Context) error
	CheckDevice(ctx context.Context) (*model.SystemMetrics, error)
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

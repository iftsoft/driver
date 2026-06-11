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

type DeviceCreator interface {
	CreateDevice(log *slog.Logger, cfg *config.DeviceConfig, cb Callback) (any, error)
}

type DummyEngine struct {
	BaseEngine
}

func NewDummyEngine(log *slog.Logger, cfg *config.DeviceConfig, cb Callback) *DummyEngine {
	return &DummyEngine{
		BaseEngine: BaseEngine{
			Log:      log,
			Config:   cfg,
			Callback: cb,
		},
	}
}

type DummyCreator struct{}

func (dc *DummyCreator) CreateDevice(log *slog.Logger, cfg *config.DeviceConfig, cb Callback) (any, error) {
	dummy := NewDummyEngine(log, cfg, cb)
	return dummy, nil
}

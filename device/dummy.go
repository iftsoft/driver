package device

import (
	"context"

	"github.com/iftsoft/linker/model"
)

type DummyCreator struct{}

func (dc *DummyCreator) CreateDevice(params CreatorParams) (any, error) {
	dummy := NewDummyEngine(params)
	return dummy, nil
}

type DummyEngine struct {
	BaseEngine
}

func NewDummyEngine(params CreatorParams) *DummyEngine {
	return &DummyEngine{
		BaseEngine: BaseEngine{
			DevName:  params.DevName,
			Log:      params.Logger,
			Config:   params.Config,
			Callback: params.Callback,
		},
	}
}

func (dm *DummyEngine) InitDevice(ctx context.Context) error {
	return nil
}

func (dm *DummyEngine) StartDevice(ctx context.Context, query *model.SystemConfig) error {
	return nil
}

func (dm *DummyEngine) StopDevice(ctx context.Context) error {
	return nil
}

func (dm *DummyEngine) CheckDevice(ctx context.Context) (*model.SystemMetrics, error) {
	out := &model.SystemMetrics{}
	return out, nil
}

func (dm *DummyEngine) DeviceTimer(ctx context.Context, unix int64) error {
	return nil
}

package device

import (
	"context"
	"log/slog"

	"github.com/iftsoft/linker/model"

	"github.com/iftsoft/driver/config"
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

func (dm *DummyEngine) DeviceSettings() model.DeviceSetup {
	return model.DeviceSetup{
		DevType:     model.DevTypeCustom,
		Supported:   model.ScopeFlagSystem | model.ScopeFlagDevice,
		Required:    model.ScopeFlagSystem | model.ScopeFlagDevice,
		Description: "Dummy device greeting",
	}
}

func (dm *DummyEngine) InitDevice(ctx context.Context) error {
	dm.Log.Info("Dummy: engine init")
	return nil
}

func (dm *DummyEngine) StartDevice(ctx context.Context, query *model.ConfigUpdate) error {
	if dm.Config.Linker == nil {
		dm.Config.Linker = &config.LinkerConfig{}
	}
	cfg := dm.Config.Linker
	cfg.LinkType = config.EnumLinkType(query.LinkType)
	cfg.Serial.PortName = query.PortName
	cfg.HidUsb.VendorID = query.VendorID
	cfg.HidUsb.ProductID = query.ProductID
	cfg.HidUsb.SerialNo = query.SerialNo
	dm.Log.Info("Dummy: device start")
	return nil
}

func (dm *DummyEngine) StopDevice(ctx context.Context) error {
	dm.Log.Info("Dummy: device stop")
	return nil
}

func (dm *DummyEngine) CheckDevice(ctx context.Context) (*model.DeviceMetrics, error) {
	out := &model.DeviceMetrics{}
	dm.Log.Info("Dummy: device check")
	return out, nil
}

func (dm *DummyEngine) DeviceTimer(ctx context.Context, unix int64) error {
	dm.Log.Info("Dummy: device timer")
	return nil
}

// Cancel interrupts current operation on device
func (dm *DummyEngine) Cancel(ctx context.Context, query *model.DeviceQuery) (*model.DeviceReply, error) {
	// TODO: Add some useful action here
	reply := dm.GetDeviceReply(model.CmdDeviceCancel)
	dm.Log.Warn("DummyDevice.Cancel", slog.Any("query", *query), slog.Any("reply", reply))
	return &reply, nil
}

// Reset initializes device to initial state
func (dm *DummyEngine) Reset(ctx context.Context, query *model.DeviceQuery) (*model.DeviceReply, error) {
	// TODO: Add some useful action here
	reply := dm.GetDeviceReply(model.CmdDeviceReset)
	dm.Log.Warn("DummyDevice.Reset", slog.Any("query", *query), slog.Any("reply", reply))
	return &reply, nil
}

// Status returns current status of device
func (dm *DummyEngine) Status(ctx context.Context, query *model.DeviceQuery) (*model.DeviceReply, error) {
	// TODO: Add some useful action here
	reply := dm.GetDeviceReply(model.CmdDeviceStatus)
	dm.Log.Warn("DummyDevice.Status", slog.Any("query", *query), slog.Any("reply", reply))
	return &reply, nil
}

// Execute returns result of command execution
func (dm *DummyEngine) Execute(ctx context.Context, query *model.DeviceQuery) (*model.DeviceReply, error) {
	// TODO: Add some useful action here
	reply := dm.GetDeviceReply(model.CmdDeviceExecute)
	dm.Log.Warn("DummyDevice.Execute", slog.Any("query", *query), slog.Any("reply", reply))
	return &reply, nil
}

package system

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/iftsoft/linker/grpc/server"
	"github.com/iftsoft/linker/model"

	"github.com/iftsoft/driver/config"
	"github.com/iftsoft/driver/utils"
)

type DeviceDriver interface {
	InitDevice(ctx context.Context) error
	StartDevice(ctx context.Context, query *model.SystemConfig) error
	StopDevice(ctx context.Context) error
	DeviceTimer(ctx context.Context, unix int64) error
	CheckDevice(ctx context.Context, metrics *model.SystemMetrics) error
}

type CallbackAPI interface {
	model.SystemCallback
	model.DeviceCallback
}

var ErrNotImplemented = errors.New("method is not implemented")

type SystemDriver struct {
	devName  string
	log      *slog.Logger
	config   *config.DeviceConfig
	grpcSrv  *server.Server
	callback CallbackAPI
	//	system   *model.SystemManager
	device model.DeviceManager
	driver DeviceDriver
	state  model.SysState
	error  model.SysError

	done    chan struct{}
	wg      sync.WaitGroup
	checkTm time.Time
}

func NewSystemDriver(log *slog.Logger, cfg *config.AppConfig, cb CallbackAPI) *SystemDriver {
	if cfg == nil {
		return nil
	}
	sd := SystemDriver{
		//devName:   cfg.DevName,
		log:      log,
		state:    model.SysStateUndefined,
		error:    model.SysErrSuccess,
		driver:   nil,
		config:   cfg.Device,
		callback: cb,
		done:     make(chan struct{}),
	}
	return &sd
}

// Implementation of model.SystemManager
func (sd *SystemDriver) Terminate(ctx context.Context, query *model.SystemQuery) (*model.SystemReply, error) {
	sd.state = model.SysStateUndefined
	sd.error = model.SysErrSuccess
	var err error
	if sd.driver != nil {
		err = sd.driver.StopDevice(ctx)
	}
	reply := &model.SystemReply{
		Device:   sd.devName,
		Command:  model.CmdSystemTerminate,
		SysState: sd.state,
		SysError: sd.error,
	}
	if err != nil {
		reply.SysError = model.SysErrSystemFail
		reply.Message = err.Error()
	}
	utils.SendQuitSignal(100)
	return reply, nil
}

func (sd *SystemDriver) SysInform(ctx context.Context, query *model.SystemQuery) (*model.SystemHealth, error) {
	reply := &model.SystemHealth{
		Device:   sd.devName,
		Moment:   time.Now().Unix(),
		SysState: sd.state,
		SysError: sd.error,
	}
	var err error
	if sd.driver != nil {
		err = sd.driver.CheckDevice(ctx, &reply.Metrics)
	}
	if err != nil {
		reply.SysError = model.SysErrSystemFail
	}
	return reply, nil
}

func (sd *SystemDriver) SysStart(ctx context.Context, query *model.SystemConfig) (*model.SystemReply, error) {
	sd.state = model.SysStateUndefined
	var err error
	if sd.driver != nil {
		err = sd.driver.StartDevice(ctx, query)
		if err == nil {
			sd.state = model.SysStateRunning
		} else {
			sd.state = model.SysStateFailed
		}
	}
	reply := &model.SystemReply{
		Device:   sd.devName,
		Command:  model.CmdSystemStart,
		SysState: sd.state,
		SysError: sd.error,
	}
	if err != nil {
		reply.SysError = model.SysErrSystemFail
		reply.Message = err.Error()
	}
	return reply, nil
}

func (sd *SystemDriver) SysStop(ctx context.Context, query *model.SystemQuery) (*model.SystemReply, error) {
	sd.state = model.SysStateUndefined
	var err error
	if sd.driver != nil {
		err = sd.driver.StopDevice(ctx)
		if err == nil {
			sd.state = model.SysStateStopped
		}
	}
	reply := &model.SystemReply{
		Device:   sd.devName,
		Command:  model.CmdSystemStop,
		SysState: sd.state,
		SysError: sd.error,
	}
	if err != nil {
		reply.Message = err.Error()
	}
	return reply, nil
}

func (sd *SystemDriver) SysRestart(ctx context.Context, query *model.SystemConfig) (*model.SystemReply, error) {
	sd.state = model.SysStateUndefined
	var err error
	if sd.driver != nil {
		err = sd.driver.StopDevice(ctx)
		if err == nil {
			sd.state = model.SysStateStopped
		}
		err = sd.driver.StartDevice(ctx, query)
		if err == nil {
			sd.state = model.SysStateRunning
		} else {
			sd.state = model.SysStateFailed
		}
	}
	reply := &model.SystemReply{
		Device:   sd.devName,
		Command:  model.CmdSystemRestart,
		SysState: sd.state,
		SysError: sd.error,
	}
	if err != nil {
		reply.Message = err.Error()
	}
	return reply, nil
}

// Cancel interrupts current operation on device
func (sd *SystemDriver) Cancel(ctx context.Context, query *model.DeviceQuery) (*model.DeviceReply, error) {
	if sd.device == nil {
		return nil, ErrNotImplemented
	}
	return sd.device.Cancel(ctx, query)
}

// Reset initializes device to initial state
func (c *SystemDriver) Reset(ctx context.Context, query *model.DeviceQuery) (*model.DeviceReply, error) {
	if c.device == nil {
		return nil, ErrNotImplemented
	}
	return c.device.Reset(ctx, query)
}

// Status returns current status of device
func (c *SystemDriver) Status(ctx context.Context, query *model.DeviceQuery) (*model.DeviceReply, error) {
	if c.device == nil {
		return nil, ErrNotImplemented
	}
	return c.device.Status(ctx, query)
}

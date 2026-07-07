package system

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/iftsoft/linker/model"

	"github.com/iftsoft/driver/utils"
)

var ErrNotImplemented = errors.New("method is not implemented")
var ErrNoDeviceDriver = errors.New("device driver is not configured")

type SystemDriver struct {
	log      *slog.Logger
	driver   *DeviceDriver
	callback model.SystemCallback
	state    model.SysState
	error    model.SysError
}

func NewSystemDriver(log *slog.Logger, callback model.SystemCallback, driver *DeviceDriver) *SystemDriver {
	sd := SystemDriver{
		log:      log,
		driver:   driver,
		callback: callback,
		state:    model.SysStateUndefined,
		error:    model.SysErrSuccess,
	}
	return &sd
}

// Implementation of model.SystemManager
func (sd *SystemDriver) Terminate(ctx context.Context, query *model.SystemQuery) (*model.SystemReply, error) {
	sd.state = model.SysStateUndefined
	sd.error = model.SysErrSuccess
	reply := &model.SystemReply{
		Device:  query.Device,
		Command: model.CmdSystemTerminate,
	}

	if sd.driver == nil {
		return reply, ErrNoDeviceDriver
	}

	err := sd.driver.DeleteDevice(ctx)
	if err != nil {
		sd.error = model.SysErrSystemFail
		reply.Message = err.Error()
	}

	reply.SysState = sd.state
	reply.SysError = sd.error
	if sd.callback != nil {
		err = sd.callback.SystemReply(ctx, reply)
	}
	go utils.SendQuitSignal(100)

	sd.log.Info("SystemDriver.Terminate", "query", *query, "reply", *reply)
	return reply, nil
}

func (sd *SystemDriver) SysInform(ctx context.Context, query *model.SystemQuery) (*model.SystemHealth, error) {
	reply := &model.SystemHealth{
		Device: query.Device,
		Moment: time.Now().Unix(),
	}

	if sd.driver == nil {
		return reply, ErrNoDeviceDriver
	}

	metrics, err := sd.driver.CheckDevice(ctx)
	if err != nil {
		sd.error = model.SysErrSystemFail
	}
	reply.SysState = sd.state
	reply.SysError = sd.error
	if metrics != nil {
		reply.Metrics = *metrics
	}
	if sd.callback != nil {
		err = sd.callback.SystemHealth(ctx, reply)
	}

	sd.log.Info("SystemDriver.SysInform", "query", *query, "reply", *reply)
	return reply, nil
}

func (sd *SystemDriver) SysStart(ctx context.Context, query *model.SystemConfig) (*model.SystemReply, error) {
	reply := &model.SystemReply{
		Device:  query.Device,
		Command: model.CmdSystemStart,
	}

	if sd.driver == nil {
		return reply, ErrNoDeviceDriver
	}

	err := sd.driver.CreateDevice(ctx, query)
	if err == nil {
		sd.state = model.SysStateRunning
	} else {
		sd.state = model.SysStateFailed
		sd.error = model.SysErrSystemFail
		reply.Message = err.Error()
	}
	reply.SysState = sd.state
	reply.SysError = sd.error
	if sd.callback != nil {
		err = sd.callback.SystemReply(ctx, reply)
	}

	sd.log.Info("SystemDriver.SysStart", "query", *query, "reply", *reply)
	return reply, nil
}

func (sd *SystemDriver) SysStop(ctx context.Context, query *model.SystemQuery) (*model.SystemReply, error) {
	reply := &model.SystemReply{
		Device:  query.Device,
		Command: model.CmdSystemStop,
	}

	if sd.driver == nil {
		return reply, ErrNoDeviceDriver
	}

	err := sd.driver.DeleteDevice(ctx)
	if err == nil {
		sd.state = model.SysStateStopped
	} else {
		sd.state = model.SysStateUndefined
		sd.error = model.SysErrSystemFail
		reply.Message = err.Error()
	}
	reply.SysState = sd.state
	reply.SysError = sd.error
	if sd.callback != nil {
		err = sd.callback.SystemReply(ctx, reply)
	}

	sd.log.Info("SystemDriver.SysStop", "query", *query, "reply", *reply)
	return reply, nil
}

func (sd *SystemDriver) SysRestart(ctx context.Context, query *model.SystemConfig) (*model.SystemReply, error) {
	sd.state = model.SysStateUndefined
	reply := &model.SystemReply{
		Device:  query.Device,
		Command: model.CmdSystemRestart,
	}

	if sd.driver == nil {
		return reply, ErrNoDeviceDriver
	}

	err := sd.driver.DeleteDevice(ctx)
	if err == nil {
		sd.state = model.SysStateStopped
	}
	err = sd.driver.CreateDevice(ctx, query)
	if err == nil {
		sd.state = model.SysStateRunning
	} else {
		sd.state = model.SysStateFailed
		reply.Message = err.Error()
	}
	reply.SysState = sd.state
	reply.SysError = sd.error
	if sd.callback != nil {
		err = sd.callback.SystemReply(ctx, reply)
	}

	sd.log.Info("SystemDriver.SysRestart", "query", *query, "reply", *reply)
	return reply, nil
}

// Cancel interrupts current operation on device
func (sd *SystemDriver) Cancel(ctx context.Context, query *model.DeviceQuery) (*model.DeviceReply, error) {
	if sd.driver == nil {
		return nil, ErrNoDeviceDriver
	}
	manager := sd.driver.DeviceManager()
	return manager.Cancel(ctx, query)
}

// Reset initializes device to initial state
func (sd *SystemDriver) Reset(ctx context.Context, query *model.DeviceQuery) (*model.DeviceReply, error) {
	if sd.driver == nil {
		return nil, ErrNoDeviceDriver
	}
	manager := sd.driver.DeviceManager()
	return manager.Reset(ctx, query)
}

// Status returns current status of device
func (sd *SystemDriver) Status(ctx context.Context, query *model.DeviceQuery) (*model.DeviceReply, error) {
	if sd.driver == nil {
		return nil, ErrNoDeviceDriver
	}
	manager := sd.driver.DeviceManager()
	return manager.Status(ctx, query)
}

// Execute returns result of command execution
func (sd *SystemDriver) Execute(ctx context.Context, query *model.DeviceQuery) (*model.DeviceReply, error) {
	if sd.driver == nil {
		return nil, ErrNoDeviceDriver
	}
	manager := sd.driver.DeviceManager()
	return manager.Status(ctx, query)
}

// InitPrinter does primary initialization of printer before printing
func (sd *SystemDriver) InitPrinter(ctx context.Context, query *model.PrinterSetup) (*model.DeviceReply, error) {
	if sd.driver == nil {
		return nil, ErrNoDeviceDriver
	}
	manager := sd.driver.PrinterManager()
	return manager.InitPrinter(ctx, query)
}

// PrintPage trys to print given text on the printer
func (sd *SystemDriver) PrintPage(ctx context.Context, query *model.PrinterQuery) (*model.DeviceReply, error) {
	if sd.driver == nil {
		return nil, ErrNoDeviceDriver
	}
	manager := sd.driver.PrinterManager()
	return manager.PrintPage(ctx, query)
}

// EnterCard trys to accept card in card reader device
func (sd *SystemDriver) EnterCard(ctx context.Context, query *model.DeviceQuery) (*model.DeviceReply, error) {
	if sd.driver == nil {
		return nil, ErrNoDeviceDriver
	}
	manager := sd.driver.ReaderManager()
	return manager.EnterCard(ctx, query)
}

// EjectCard trys to reject card from card reader device
func (sd *SystemDriver) EjectCard(ctx context.Context, query *model.DeviceQuery) (*model.DeviceReply, error) {
	if sd.driver == nil {
		return nil, ErrNoDeviceDriver
	}
	manager := sd.driver.ReaderManager()
	return manager.EjectCard(ctx, query)
}

// CaptureCard trys to capture card inside card reader device
func (sd *SystemDriver) CaptureCard(ctx context.Context, query *model.DeviceQuery) (*model.DeviceReply, error) {
	if sd.driver == nil {
		return nil, ErrNoDeviceDriver
	}
	manager := sd.driver.ReaderManager()
	return manager.CaptureCard(ctx, query)
}

// ReadCard trys to read card information from card
func (sd *SystemDriver) ReadCard(ctx context.Context, query *model.DeviceQuery) (*model.ReadCardReply, error) {
	if sd.driver == nil {
		return nil, ErrNoDeviceDriver
	}
	manager := sd.driver.ReaderManager()
	return manager.ReadCard(ctx, query)
}

// InitValidator does primary initialization of the validator
func (sd *SystemDriver) InitValidator(ctx context.Context, query *model.ValidatorQuery) (*model.DeviceReply, error) {
	if sd.driver == nil {
		return nil, ErrNoDeviceDriver
	}
	manager := sd.driver.ValidatorManager()
	return manager.InitValidator(ctx, query)
}

// DoValidate starts accepting cash from user
func (sd *SystemDriver) DoValidate(ctx context.Context, query *model.ValidatorQuery) (*model.DeviceReply, error) {
	if sd.driver == nil {
		return nil, ErrNoDeviceDriver
	}
	manager := sd.driver.ValidatorManager()
	return manager.DoValidate(ctx, query)
}

// AcceptNote puts the validated note to the cassette
func (sd *SystemDriver) AcceptNote(ctx context.Context, query *model.ValidatorQuery) (*model.DeviceReply, error) {
	if sd.driver == nil {
		return nil, ErrNoDeviceDriver
	}
	manager := sd.driver.ValidatorManager()
	return manager.AcceptNote(ctx, query)
}

// ReturnNote returns the validated note to the user
func (sd *SystemDriver) ReturnNote(ctx context.Context, query *model.ValidatorQuery) (*model.DeviceReply, error) {
	if sd.driver == nil {
		return nil, ErrNoDeviceDriver
	}
	manager := sd.driver.ValidatorManager()
	return manager.ReturnNote(ctx, query)
}

// StopValidate disables accepting new notes by validator
func (sd *SystemDriver) StopValidate(ctx context.Context, query *model.ValidatorQuery) (*model.DeviceReply, error) {
	if sd.driver == nil {
		return nil, ErrNoDeviceDriver
	}
	manager := sd.driver.ValidatorManager()
	return manager.StopValidate(ctx, query)
}

// CheckValidator returns current cassette state
func (sd *SystemDriver) CheckValidator(ctx context.Context, query *model.ValidatorQuery) (*model.ValidatorStore, error) {
	if sd.driver == nil {
		return nil, ErrNoDeviceDriver
	}
	manager := sd.driver.ValidatorManager()
	return manager.CheckValidator(ctx, query)
}

// ClearValidator clears all cassette data (settlement or reconciliation)
func (sd *SystemDriver) ClearValidator(ctx context.Context, query *model.ValidatorQuery) (*model.ValidatorStore, error) {
	if sd.driver == nil {
		return nil, ErrNoDeviceDriver
	}
	manager := sd.driver.ValidatorManager()
	return manager.ClearValidator(ctx, query)
}

package system

import (
	"context"
	"log/slog"

	"github.com/iftsoft/linker/model"
)

type NoopDevice struct {
	log *slog.Logger
}

func NewNoopDevice(log *slog.Logger) *NoopDevice {
	return &NoopDevice{log: log}
}

func (dum *NoopDevice) DeviceSettings() model.DeviceSetup {
	return model.DeviceSetup{}
}

func (dum *NoopDevice) InitDevice(ctx context.Context) error {
	return ErrNotImplemented
}

func (dum *NoopDevice) StartDevice(ctx context.Context, query *model.ConfigUpdate) error {
	return ErrNotImplemented
}

func (dum *NoopDevice) StopDevice(ctx context.Context) error {
	return nil
}

func (dum *NoopDevice) CheckDevice(ctx context.Context) (*model.DeviceMetrics, error) {
	return nil, ErrNotImplemented
}

func (dum *NoopDevice) DeviceTimer(ctx context.Context, unix int64) error {
	return ErrNotImplemented
}

// Cancel interrupts current operation on device
func (dum *NoopDevice) Cancel(ctx context.Context, query *model.DeviceQuery) (*model.DeviceReply, error) {
	if query != nil {
		dum.log.Warn("NoopDevice.Cancel", "query", *query)
	}
	return nil, ErrNotImplemented
}

// Reset initializes device to initial state
func (dum *NoopDevice) Reset(ctx context.Context, query *model.DeviceQuery) (*model.DeviceReply, error) {
	if query != nil {
		dum.log.Warn("NoopDevice.Reset", "query", *query)
	}
	return nil, ErrNotImplemented
}

// Status returns current status of device
func (dum *NoopDevice) Status(ctx context.Context, query *model.DeviceQuery) (*model.DeviceReply, error) {
	if query != nil {
		dum.log.Warn("NoopDevice.Status", "query", *query)
	}
	return nil, ErrNotImplemented
}

// Execute returns result of command execution
func (dum *NoopDevice) Execute(ctx context.Context, query *model.DeviceQuery) (*model.DeviceReply, error) {
	if query != nil {
		dum.log.Warn("NoopDevice.Execute", "query", *query)
	}
	return nil, ErrNotImplemented
}

// InitPrinter does primary initialization of printer before printing
func (dum *NoopDevice) InitPrinter(ctx context.Context, query *model.PrinterSetup) (*model.DeviceReply, error) {
	if query != nil {
		dum.log.Warn("NoopDevice.InitPrinter", "query", *query)
	}
	return nil, ErrNotImplemented
}

// PrintPage trys to print given text on the printer
func (dum *NoopDevice) PrintPage(ctx context.Context, query *model.PrinterQuery) (*model.DeviceReply, error) {
	if query != nil {
		dum.log.Warn("NoopDevice.PrintPage", "query", *query)
	}
	return nil, ErrNotImplemented
}

// EnterCard trys to accept card in card reader device
func (dum *NoopDevice) EnterCard(ctx context.Context, query *model.DeviceQuery) (*model.DeviceReply, error) {
	if query != nil {
		dum.log.Warn("NoopDevice.EnterCard", "query", *query)
	}
	return nil, ErrNotImplemented
}

// EjectCard trys to reject card from card reader device
func (dum *NoopDevice) EjectCard(ctx context.Context, query *model.DeviceQuery) (*model.DeviceReply, error) {
	if query != nil {
		dum.log.Warn("NoopDevice.EjectCard", "query", *query)
	}
	return nil, ErrNotImplemented
}

// CaptureCard trys to capture card inside card reader device
func (dum *NoopDevice) CaptureCard(ctx context.Context, query *model.DeviceQuery) (*model.DeviceReply, error) {
	if query != nil {
		dum.log.Warn("NoopDevice.CaptureCard", "query", *query)
	}
	return nil, ErrNotImplemented
}

// ReadCard trys to read card information from card
func (dum *NoopDevice) ReadCard(ctx context.Context, query *model.DeviceQuery) (*model.ReadCardReply, error) {
	if query != nil {
		dum.log.Warn("NoopDevice.ReadCard", "query", *query)
	}
	return nil, ErrNotImplemented
}

// InitValidator does primary initialization of the validator
func (dum *NoopDevice) InitValidator(ctx context.Context, query *model.ValidatorQuery) (*model.DeviceReply, error) {
	if query != nil {
		dum.log.Warn("NoopDevice.InitValidator", "query", *query)
	}
	return nil, ErrNotImplemented
}

// DoValidate starts accepting cash from user
func (dum *NoopDevice) DoValidate(ctx context.Context, query *model.ValidatorQuery) (*model.DeviceReply, error) {
	if query != nil {
		dum.log.Warn("NoopDevice.DoValidate", "query", *query)
	}
	return nil, ErrNotImplemented
}

// AcceptNote puts the validated note to the cassette
func (dum *NoopDevice) AcceptNote(ctx context.Context, query *model.ValidatorQuery) (*model.DeviceReply, error) {
	if query != nil {
		dum.log.Warn("NoopDevice.AcceptNote", "query", *query)
	}
	return nil, ErrNotImplemented
}

// ReturnNote returns the validated note to the user
func (dum *NoopDevice) ReturnNote(ctx context.Context, query *model.ValidatorQuery) (*model.DeviceReply, error) {
	if query != nil {
		dum.log.Warn("NoopDevice.ReturnNote", "query", *query)
	}
	return nil, ErrNotImplemented
}

// StopValidate disables accepting new notes by validator
func (dum *NoopDevice) StopValidate(ctx context.Context, query *model.ValidatorQuery) (*model.DeviceReply, error) {
	if query != nil {
		dum.log.Warn("NoopDevice.StopValidate", "query", *query)
	}
	return nil, ErrNotImplemented
}

// CheckValidator returns current cassette state
func (dum *NoopDevice) CheckValidator(ctx context.Context, query *model.ValidatorQuery) (*model.ValidatorStore, error) {
	if query != nil {
		dum.log.Warn("NoopDevice.CheckValidator", "query", *query)
	}
	return nil, ErrNotImplemented
}

// ClearValidator clears all cassette data (settlement or reconciliation)
func (dum *NoopDevice) ClearValidator(ctx context.Context, query *model.ValidatorQuery) (*model.ValidatorStore, error) {
	if query != nil {
		dum.log.Warn("NoopDevice.ClearValidator", "query", *query)
	}
	return nil, ErrNotImplemented
}

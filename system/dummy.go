package system

import (
	"context"
	"log/slog"

	"github.com/iftsoft/linker/model"
)

type DummyDevice struct {
	log *slog.Logger
}

func NewDummyDevice(log *slog.Logger) *DummyDevice {
	return &DummyDevice{log: log}
}

func (dum *DummyDevice) DeviceSettings() model.DeviceSetup {
	return model.DeviceSetup{}
}

func (dum *DummyDevice) InitDevice(ctx context.Context) error {
	return ErrNotImplemented
}

func (dum *DummyDevice) StartDevice(ctx context.Context, query *model.ConfigUpdate) error {
	return ErrNotImplemented
}

func (dum *DummyDevice) StopDevice(ctx context.Context) error {
	return ErrNotImplemented
}

func (dum *DummyDevice) CheckDevice(ctx context.Context) (*model.DeviceMetrics, error) {
	return nil, ErrNotImplemented
}

func (dum *DummyDevice) DeviceTimer(ctx context.Context, unix int64) error {
	return ErrNotImplemented
}

// Cancel interrupts current operation on device
func (dum *DummyDevice) Cancel(ctx context.Context, query *model.DeviceQuery) (*model.DeviceReply, error) {
	dum.log.Warn("DummyDevice.Cancel", "query", *query)
	return nil, ErrNotImplemented
}

// Reset initializes device to initial state
func (dum *DummyDevice) Reset(ctx context.Context, query *model.DeviceQuery) (*model.DeviceReply, error) {
	dum.log.Warn("DummyDevice.Reset", "query", *query)
	return nil, ErrNotImplemented
}

// Status returns current status of device
func (dum *DummyDevice) Status(ctx context.Context, query *model.DeviceQuery) (*model.DeviceReply, error) {
	dum.log.Warn("DummyDevice.Status", "query", *query)
	return nil, ErrNotImplemented
}

// Execute returns result of command execution
func (dum *DummyDevice) Execute(ctx context.Context, query *model.DeviceQuery) (*model.DeviceReply, error) {
	dum.log.Warn("DummyDevice.Execute", "query", *query)
	return nil, ErrNotImplemented
}

// InitPrinter does primary initialization of printer before printing
func (dum *DummyDevice) InitPrinter(ctx context.Context, query *model.PrinterSetup) (*model.DeviceReply, error) {
	dum.log.Warn("DummyDevice.InitPrinter", "query", *query)
	return nil, ErrNotImplemented
}

// PrintPage trys to print given text on the printer
func (dum *DummyDevice) PrintPage(ctx context.Context, query *model.PrinterQuery) (*model.DeviceReply, error) {
	dum.log.Warn("DummyDevice.PrintPage", "query", *query)
	return nil, ErrNotImplemented
}

// EnterCard trys to accept card in card reader device
func (dum *DummyDevice) EnterCard(ctx context.Context, query *model.DeviceQuery) (*model.DeviceReply, error) {
	dum.log.Warn("DummyDevice.EnterCard", "query", *query)
	return nil, ErrNotImplemented
}

// EjectCard trys to reject card from card reader device
func (dum *DummyDevice) EjectCard(ctx context.Context, query *model.DeviceQuery) (*model.DeviceReply, error) {
	dum.log.Warn("DummyDevice.EjectCard", "query", *query)
	return nil, ErrNotImplemented
}

// CaptureCard trys to capture card inside card reader device
func (dum *DummyDevice) CaptureCard(ctx context.Context, query *model.DeviceQuery) (*model.DeviceReply, error) {
	dum.log.Warn("DummyDevice.CaptureCard", "query", *query)
	return nil, ErrNotImplemented
}

// ReadCard trys to read card information from card
func (dum *DummyDevice) ReadCard(ctx context.Context, query *model.DeviceQuery) (*model.ReadCardReply, error) {
	dum.log.Warn("DummyDevice.ReadCard", "query", *query)
	return nil, ErrNotImplemented
}

// InitValidator does primary initialization of the validator
func (dum *DummyDevice) InitValidator(ctx context.Context, query *model.ValidatorQuery) (*model.DeviceReply, error) {
	dum.log.Warn("DummyDevice.InitValidator", "query", *query)
	return nil, ErrNotImplemented
}

// DoValidate starts accepting cash from user
func (dum *DummyDevice) DoValidate(ctx context.Context, query *model.ValidatorQuery) (*model.DeviceReply, error) {
	dum.log.Warn("DummyDevice.DoValidate", "query", *query)
	return nil, ErrNotImplemented
}

// AcceptNote puts the validated note to the cassette
func (dum *DummyDevice) AcceptNote(ctx context.Context, query *model.ValidatorQuery) (*model.DeviceReply, error) {
	dum.log.Warn("DummyDevice.AcceptNote", "query", *query)
	return nil, ErrNotImplemented
}

// ReturnNote returns the validated note to the user
func (dum *DummyDevice) ReturnNote(ctx context.Context, query *model.ValidatorQuery) (*model.DeviceReply, error) {
	dum.log.Warn("DummyDevice.ReturnNote", "query", *query)
	return nil, ErrNotImplemented
}

// StopValidate disables accepting new notes by validator
func (dum *DummyDevice) StopValidate(ctx context.Context, query *model.ValidatorQuery) (*model.DeviceReply, error) {
	dum.log.Warn("DummyDevice.StopValidate", "query", *query)
	return nil, ErrNotImplemented
}

// CheckValidator returns current cassette state
func (dum *DummyDevice) CheckValidator(ctx context.Context, query *model.ValidatorQuery) (*model.ValidatorStore, error) {
	dum.log.Warn("DummyDevice.CheckValidator", "query", *query)
	return nil, ErrNotImplemented
}

// ClearValidator clears all cassette data (settlement or reconciliation)
func (dum *DummyDevice) ClearValidator(ctx context.Context, query *model.ValidatorQuery) (*model.ValidatorStore, error) {
	dum.log.Warn("DummyDevice.ClearValidator", "query", *query)
	return nil, ErrNotImplemented
}

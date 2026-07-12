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
		BaseEngine: NewBaseEngine(params),
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

//func (dm *DummyEngine) InitDevice(ctx context.Context) error {
//	dm.Log.Info("Dummy: engine init")
//	return nil
//}

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
	dm.Log.Debug("Dummy: device start")
	// TODO: run some initialization here
	dm.DevState = model.DevStateReady
	return nil
}

func (dm *DummyEngine) StopDevice(ctx context.Context) error {
	dm.Log.Debug("Dummy: device stop")
	// TODO: run some cleanup here
	dm.ClearDevice()
	return nil
}

func (dm *DummyEngine) CheckDevice(ctx context.Context) (*model.DeviceMetrics, error) {
	dm.Log.Debug("Dummy: device check")
	out := dm.GetDeviceMetrics()
	// TODO: add extra content here
	return out, nil
}

func (dm *DummyEngine) DeviceTimer(ctx context.Context, unix int64) error {
	dm.Log.Debug("Dummy: device timer", slog.Int64("unix", unix))
	return nil
}

// Cancel interrupts current operation on device
func (dm *DummyEngine) Cancel(ctx context.Context, query *model.DeviceQuery) (*model.DeviceReply, error) {
	// TODO: Add some useful action here
	reply := dm.GetDeviceReply(model.CmdDeviceCancel)
	dm.NotifyDeviceReply(ctx, reply)
	dm.Log.Warn("DummyDevice.Cancel", slog.Any("query", query), slog.Any("reply", reply))
	return &reply, nil
}

// Reset initializes device to initial state
func (dm *DummyEngine) Reset(ctx context.Context, query *model.DeviceQuery) (*model.DeviceReply, error) {
	dm.ClearDevice()
	// TODO: Add some useful action here
	reply := dm.GetDeviceReply(model.CmdDeviceReset)
	dm.NotifyDeviceReply(ctx, reply)
	dm.Log.Warn("DummyDevice.Reset", slog.Any("query", query), slog.Any("reply", reply))
	return &reply, nil
}

// Status returns current status of device
func (dm *DummyEngine) Status(ctx context.Context, query *model.DeviceQuery) (*model.DeviceReply, error) {
	// TODO: Add some useful action here
	reply := dm.GetDeviceReply(model.CmdDeviceStatus)
	dm.NotifyDeviceReply(ctx, reply)
	dm.Log.Warn("DummyDevice.Status", slog.Any("query", query), slog.Any("reply", reply))
	return &reply, nil
}

// Execute returns result of command execution
func (dm *DummyEngine) Execute(ctx context.Context, query *model.DeviceQuery) (*model.DeviceReply, error) {
	dm.ClearCommand()
	// TODO: Add some useful action here
	dm.DevResult = "12345678"
	dm.NotifyReaderResult(ctx, dm.DevResult)
	reply := dm.GetDeviceReply(model.CmdDeviceExecute)
	dm.NotifyDeviceReply(ctx, reply)
	dm.Log.Warn("DummyDevice.Execute", slog.Any("query", query), slog.Any("reply", reply))
	return &reply, nil
}

// InitPrinter does primary initialization of printer before printing
func (dm *DummyEngine) InitPrinter(ctx context.Context, query *model.PrinterSetup) (*model.DeviceReply, error) {
	dm.ClearCommand()
	// TODO: Add some useful action here
	reply := dm.GetDeviceReply(model.CmdInitPrinter)
	dm.NotifyDeviceReply(ctx, reply)
	dm.Log.Warn("DummyDevice.InitPrinter", slog.Any("query", query), slog.Any("reply", reply))
	return &reply, nil
}

// PrintPage trys to print given text on the printer
func (dm *DummyEngine) PrintPage(ctx context.Context, query *model.PrinterQuery) (*model.DeviceReply, error) {
	dm.ClearCommand()
	// TODO: Add some useful action here
	reply := dm.GetDeviceReply(model.CmdPrintPage)
	dm.NotifyDeviceReply(ctx, reply)
	dm.Log.Warn("DummyDevice.PrintPage", slog.Any("query", query), slog.Any("reply", reply))
	return &reply, nil
}

// EnterCard trys to accept card in card reader device
func (dm *DummyEngine) EnterCard(ctx context.Context, query *model.DeviceQuery) (*model.DeviceReply, error) {
	dm.ClearCommand()
	// TODO: Add some useful action here
	position := model.PositionNotify{Position: 1}
	dm.NotifyCardPosition(ctx, position)
	reply := dm.GetDeviceReply(model.CmdEnterCard)
	dm.NotifyDeviceReply(ctx, reply)
	dm.Log.Warn("DummyDevice.EnterCard", slog.Any("query", query), slog.Any("reply", reply))
	return &reply, nil
}

// EjectCard trys to reject card from card reader device
func (dm *DummyEngine) EjectCard(ctx context.Context, query *model.DeviceQuery) (*model.DeviceReply, error) {
	dm.ClearCommand()
	// TODO: Add some useful action here
	position := model.PositionNotify{Position: 0}
	dm.NotifyCardPosition(ctx, position)
	reply := dm.GetDeviceReply(model.CmdEjectCard)
	dm.NotifyDeviceReply(ctx, reply)
	dm.Log.Warn("DummyDevice.EjectCard", slog.Any("query", query), slog.Any("reply", reply))
	return &reply, nil
}

// CaptureCard trys to capture card inside card reader device
func (dm *DummyEngine) CaptureCard(ctx context.Context, query *model.DeviceQuery) (*model.DeviceReply, error) {
	dm.ClearCommand()
	// TODO: Add some useful action here
	position := model.PositionNotify{Position: 2}
	dm.NotifyCardPosition(ctx, position)
	reply := dm.GetDeviceReply(model.CmdCaptureCard)
	dm.NotifyDeviceReply(ctx, reply)
	dm.Log.Warn("DummyDevice.CaptureCard", slog.Any("query", query), slog.Any("reply", reply))
	return &reply, nil
}

// ReadCard trys to read card information from card
func (dm *DummyEngine) ReadCard(ctx context.Context, query *model.DeviceQuery) (*model.ReadCardReply, error) {
	dm.ClearCommand()
	// TODO: Add some useful action here
	content := model.CardContent{
		CardPan: "411111000012345678",
		ExpDate: "10/30",
		Holder:  "John Doe",
	}
	dm.NotifyCardDescription(ctx, content)
	reply := model.ReadCardReply{
		DeviceReply: dm.GetDeviceReply(model.CmdReadCard),
		CardContent: content,
	}
	dm.NotifyDeviceReply(ctx, reply.DeviceReply)
	dm.Log.Warn("DummyDevice.ReadCard", slog.Any("query", query), slog.Any("reply", reply))
	return &reply, nil
}

// InitValidator does primary initialization of the validator
func (dm *DummyEngine) InitValidator(ctx context.Context, query *model.ValidatorQuery) (*model.DeviceReply, error) {
	dm.ClearCommand()
	// TODO: Add some useful action here
	reply := dm.GetDeviceReply(model.CmdInitValidator)
	dm.NotifyDeviceReply(ctx, reply)
	dm.Log.Warn("DummyDevice.InitValidator", slog.Any("query", query), slog.Any("reply", reply))
	return &reply, nil
}

// DoValidate starts accepting cash from user
func (dm *DummyEngine) DoValidate(ctx context.Context, query *model.ValidatorQuery) (*model.DeviceReply, error) {
	dm.ClearCommand()
	// TODO: Add some useful action here
	accept := model.AcceptNotify{
		Note: model.ValidatorNote{
			Currency: query.Currency,
			Nominal:  1000,
			Count:    1,
			Amount:   1000,
		},
	}
	dm.NotifyNoteAccepted(ctx, accept)
	reply := dm.GetDeviceReply(model.CmdDoValidate)
	dm.NotifyDeviceReply(ctx, reply)
	dm.Log.Warn("DummyDevice.DoValidate", slog.Any("query", query), slog.Any("reply", reply))
	return &reply, nil
}

// AcceptNote puts the validated note to the cassette
func (dm *DummyEngine) AcceptNote(ctx context.Context, query *model.ValidatorQuery) (*model.DeviceReply, error) {
	dm.ClearCommand()
	// TODO: Add some useful action here
	accept := model.AcceptNotify{
		Note: model.ValidatorNote{
			Currency: query.Currency,
			Nominal:  1000,
			Count:    1,
			Amount:   1000,
		},
	}
	dm.NotifyCashIsStored(ctx, accept)
	reply := dm.GetDeviceReply(model.CmdAcceptNote)
	dm.NotifyDeviceReply(ctx, reply)
	dm.Log.Warn("DummyDevice.AcceptNote", slog.Any("query", query), slog.Any("reply", reply))
	return &reply, nil
}

// ReturnNote returns the validated note to the user
func (dm *DummyEngine) ReturnNote(ctx context.Context, query *model.ValidatorQuery) (*model.DeviceReply, error) {
	dm.ClearCommand()
	// TODO: Add some useful action here
	accept := model.AcceptNotify{
		Note: model.ValidatorNote{
			Currency: query.Currency,
			Nominal:  1000,
			Count:    1,
			Amount:   1000,
		},
	}
	dm.NotifyCashReturned(ctx, accept)
	reply := dm.GetDeviceReply(model.CmdReturnNote)
	dm.NotifyDeviceReply(ctx, reply)
	dm.Log.Warn("DummyDevice.ReturnNote", slog.Any("query", query), slog.Any("reply", reply))
	return &reply, nil
}

// StopValidate disables accepting new notes by validator
func (dm *DummyEngine) StopValidate(ctx context.Context, query *model.ValidatorQuery) (*model.DeviceReply, error) {
	dm.ClearCommand()
	// TODO: Add some useful action here
	reply := dm.GetDeviceReply(model.CmdStopValidate)
	dm.NotifyDeviceReply(ctx, reply)
	dm.Log.Warn("DummyDevice.StopValidate", slog.Any("query", query), slog.Any("reply", reply))
	return &reply, nil
}

// CheckValidator returns current cassette state
func (dm *DummyEngine) CheckValidator(ctx context.Context, query *model.ValidatorQuery) (*model.ValidatorStore, error) {
	dm.ClearCommand()
	// TODO: Add some useful action here
	content := model.BatchContent{
		BatchId:    123,
		BatchState: model.StateActive,
		Notes: []model.ValidatorNote{
			{
				Currency: query.Currency,
				Nominal:  1000,
				Count:    2,
				Amount:   2000,
			},
			{
				Currency: query.Currency,
				Nominal:  2000,
				Count:    1,
				Amount:   2000,
			},
		},
	}
	dm.NotifyValidatorStore(ctx, content)
	reply := model.ValidatorStore{
		DeviceReply:  dm.GetDeviceReply(model.CmdCheckValidator),
		BatchContent: content,
	}
	dm.NotifyDeviceReply(ctx, reply.DeviceReply)
	dm.Log.Warn("DummyDevice.CheckValidator", slog.Any("query", query), slog.Any("reply", reply))
	return &reply, nil
}

// ClearValidator clears all cassette data (settlement or reconciliation)
func (dm *DummyEngine) ClearValidator(ctx context.Context, query *model.ValidatorQuery) (*model.ValidatorStore, error) {
	dm.ClearCommand()
	// TODO: Add some useful action here
	content := model.BatchContent{
		BatchId:    123,
		BatchState: model.StateActive,
		Notes: []model.ValidatorNote{
			{
				Currency: query.Currency,
				Nominal:  1000,
				Count:    2,
				Amount:   2000,
			},
			{
				Currency: query.Currency,
				Nominal:  2000,
				Count:    1,
				Amount:   2000,
			},
		},
	}
	dm.NotifyValidatorStore(ctx, content)
	reply := model.ValidatorStore{
		DeviceReply:  dm.GetDeviceReply(model.CmdClearValidator),
		BatchContent: content,
	}
	dm.NotifyDeviceReply(ctx, reply.DeviceReply)
	dm.Log.Warn("DummyDevice.ClearValidator", slog.Any("query", query), slog.Any("reply", reply))
	return &reply, nil
}

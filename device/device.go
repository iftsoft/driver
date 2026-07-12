package device

import (
	"context"
	"log/slog"
	"time"

	"github.com/iftsoft/linker/model"

	"github.com/iftsoft/driver/config"
)

type BaseEngine struct {
	Log       *slog.Logger
	Config    *config.DeviceConfig
	Callback  Callback
	DevName   string
	DevState  model.DevState
	DevError  model.DevError
	DevPrompt model.DevPrompt
	DevAction model.DevAction
	DevResult string
	DevReply  string
	StartedAt time.Time
}

func NewBaseEngine(params CreatorParams) BaseEngine {
	return BaseEngine{
		DevName:   params.DevName,
		Log:       params.Logger,
		Config:    params.Config,
		Callback:  params.Callback,
		DevState:  model.DevStateUndefined,
		DevError:  model.DevErrorSuccess,
		DevPrompt: model.DevPromptNone,
		DevAction: model.DevActionDoNothing,
		StartedAt: time.Now(),
	}
}

func (be *BaseEngine) ClearDevice() {
	be.DevState = model.DevStateUndefined
	be.DevError = model.DevErrorSuccess
	be.ClearCommand()
}

func (be *BaseEngine) ClearCommand() {
	be.DevPrompt = model.DevPromptNone
	be.DevAction = model.DevActionDoNothing
	be.DevResult = ""
	be.DevReply = ""
}

func (be *BaseEngine) GetDeviceReply(cmd string) model.DeviceReply {
	return model.DeviceReply{
		Command: cmd,
		Device:  be.DevName,
		Action:  be.DevAction,
		State:   be.DevState,
		ErrCode: be.DevError,
		ErrText: be.DevReply,
	}
}

func (be *BaseEngine) GetDeviceNotify() model.DeviceNotify {
	return model.DeviceNotify{
		Device: be.DevName,
		Action: be.DevAction,
	}
}

func (be *BaseEngine) GetDeviceMetrics() *model.DeviceMetrics {
	now := time.Now()
	uptime := now.Sub(be.StartedAt)
	return &model.DeviceMetrics{
		Moment:   now.Unix(),
		Uptime:   uint64(uptime.Milliseconds()),
		DevError: be.DevError,
		DevState: be.DevState,
		Content:  make(map[string]any),
	}
}

// NotifyDeviceReply notifies manager about command execution results
func (be *BaseEngine) NotifyDeviceReply(ctx context.Context, reply model.DeviceReply) {
	// StateChanged processing
	var err error
	if be.Callback != nil {
		err = be.Callback.DeviceReply(ctx, &reply)
		if err != nil && be.Log != nil {
			be.Log.Warn("NotifyDeviceReply failed", slog.String("err", err.Error()))
		}
	}
}

// NotifyExecuteError explicitly notifies manager about execution error
func (be *BaseEngine) NotifyExecuteError(ctx context.Context, cmd string, errCode model.DevError, reason error) {
	be.DevError = errCode
	be.DevReply = model.ExtendError(errCode, reason).Error()
	// ExecuteError processing
	if be.DevError != model.DevErrorSuccess {
		query := be.GetDeviceReply(cmd)
		if be.Callback != nil {
			err := be.Callback.ExecuteError(ctx, &query)
			if err != nil && be.Log != nil {
				be.Log.Warn("NotifyExecuteError failed", slog.String("err", err.Error()))
			}
		}
	}
}

// NotifyStateChanged notifies manager about device state changes during command execution
func (be *BaseEngine) NotifyStateChanged(ctx context.Context, newState model.DevState) {
	// StateChanged processing
	if be.DevState != newState {
		query := model.DeviceState{
			DeviceNotify: be.GetDeviceNotify(),
			StateNotify: model.StateNotify{
				OldState: be.DevState,
				NewState: newState,
			},
		}
		if be.Callback != nil {
			err := be.Callback.StateChanged(ctx, &query)
			if err != nil && be.Log != nil {
				be.Log.Warn("NotifyStateChanged failed", slog.String("err", err.Error()))
			}
		}
		be.DevState = newState
	}
}

// NotifyActionPrompt notifies manager about a user prompt during command execution
func (be *BaseEngine) NotifyActionPrompt(ctx context.Context, prompt model.DevPrompt) {
	be.DevPrompt = prompt
	// ActionPrompt processing
	if be.DevPrompt != model.DevPromptNone {
		query := &model.DevicePrompt{
			DeviceNotify: be.GetDeviceNotify(),
			PromptNotify: model.PromptNotify{
				Prompt: be.DevPrompt,
			},
		}
		if be.Callback != nil {
			err := be.Callback.ActionPrompt(ctx, query)
			if err != nil && be.Log != nil {
				be.Log.Warn("NotifyActionPrompt failed", slog.String("err", err.Error()))
			}
		}
	}
}

// NotifyReaderResult notifies manager about a data value returned by reading operation
func (be *BaseEngine) NotifyReaderResult(ctx context.Context, result string) {
	be.DevResult = result
	// ReaderResult processing
	if be.DevResult != "" {
		query := &model.DeviceResult{
			DeviceNotify: be.GetDeviceNotify(),
			ResultNotify: model.ResultNotify{
				Result: be.DevResult,
			},
		}
		if be.Callback != nil {
			err := be.Callback.ReaderResult(ctx, query)
			if err != nil && be.Log != nil {
				be.Log.Warn("NotifyReaderResult failed", slog.String("err", err.Error()))
			}
		}
	}
}

// NotifyPrinterProgress notifies manager about progress with printing operation
func (be *BaseEngine) NotifyPrinterProgress(ctx context.Context, notify model.ProgressNotify) {
	// PrinterProgress processing
	query := &model.PrinterProgress{
		DeviceNotify:   be.GetDeviceNotify(),
		ProgressNotify: notify,
	}
	if be.Callback != nil {
		err := be.Callback.PrinterProgress(ctx, query)
		if err != nil && be.Log != nil {
			be.Log.Warn("NotifyPrinterProgress failed", slog.String("err", err.Error()))
		}
	}
}

// NotifyCardPosition notifies manager about current card position
func (be *BaseEngine) NotifyCardPosition(ctx context.Context, notify model.PositionNotify) {
	// CardPosition processing
	query := &model.CardPosition{
		DeviceNotify:   be.GetDeviceNotify(),
		PositionNotify: notify,
	}
	if be.Callback != nil {
		err := be.Callback.CardPosition(ctx, query)
		if err != nil && be.Log != nil {
			be.Log.Warn("NotifyCardPosition failed", slog.String("err", err.Error()))
		}
	}
}

// NotifyCardDescription notifies manager about content of processed card
func (be *BaseEngine) NotifyCardDescription(ctx context.Context, notify model.CardContent) {
	// CardDescription processing
	query := &model.CardDescription{
		DeviceNotify: be.GetDeviceNotify(),
		CardContent:  notify,
	}
	if be.Callback != nil {
		err := be.Callback.CardDescription(ctx, query)
		if err != nil && be.Log != nil {
			be.Log.Warn("NotifyCardDescription failed", slog.String("err", err.Error()))
		}
	}
}

// NotifyNoteAccepted notifies manager about new note in validator
func (be *BaseEngine) NotifyNoteAccepted(ctx context.Context, notify model.AcceptNotify) {
	// NoteAccepted processing
	query := &model.ValidatorAccept{
		DeviceNotify: be.GetDeviceNotify(),
		AcceptNotify: notify,
	}
	if be.Callback != nil {
		err := be.Callback.NoteAccepted(ctx, query)
		if err != nil && be.Log != nil {
			be.Log.Warn("NotifyNoteAccepted failed", slog.String("err", err.Error()))
		}
	}
}

// NotifyCashIsStored notifies manager that accepted note was stored to cassette
func (be *BaseEngine) NotifyCashIsStored(ctx context.Context, notify model.AcceptNotify) {
	// CashIsStored processing
	query := &model.ValidatorAccept{
		DeviceNotify: be.GetDeviceNotify(),
		AcceptNotify: notify,
	}
	if be.Callback != nil {
		err := be.Callback.CashIsStored(ctx, query)
		if err != nil && be.Log != nil {
			be.Log.Warn("NotifyCashIsStored failed", slog.String("err", err.Error()))
		}
	}
}

// NotifyCashReturned notifies manager that accepted note was returned to user
func (be *BaseEngine) NotifyCashReturned(ctx context.Context, notify model.AcceptNotify) {
	// CashReturned processing
	query := &model.ValidatorAccept{
		DeviceNotify: be.GetDeviceNotify(),
		AcceptNotify: notify,
	}
	if be.Callback != nil {
		err := be.Callback.CashReturned(ctx, query)
		if err != nil && be.Log != nil {
			be.Log.Warn("NotifyCashReturned failed", slog.String("err", err.Error()))
		}
	}
}

// NotifyValidatorStore notifies manager about current batch content
func (be *BaseEngine) NotifyValidatorStore(ctx context.Context, notify model.BatchContent) {
	// ValidatorStore processing
	query := &model.ValidatorBatch{
		DeviceNotify: be.GetDeviceNotify(),
		BatchContent: notify,
	}
	if be.Callback != nil {
		err := be.Callback.ValidatorStore(ctx, query)
		if err != nil && be.Log != nil {
			be.Log.Warn("NotifyValidatorStore failed", slog.String("err", err.Error()))
		}
	}
}

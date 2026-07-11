package device

import (
	"context"
	"log/slog"

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
}

func (be *BaseEngine) ClearDevice() {
	be.DevState = model.DevStateUndefined
	be.DevError = model.DevErrorSuccess
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

func (be *BaseEngine) RunDeviceReply(ctx context.Context, cmd string) error {
	// StateChanged processing
	var err error
	reply := be.GetDeviceReply(cmd)
	if be.Callback != nil {
		err = be.Callback.DeviceReply(ctx, &reply)
	}
	if be.Log != nil {
		be.Log.Debug("Callback DeviceReply", slog.Any("reply", reply))
	}
	return err
}

func (be *BaseEngine) RunExecuteError(ctx context.Context, cmd string, errCode model.DevError, reason string) error {
	be.DevError = errCode
	be.DevReply = model.NewError(errCode, reason).Error()
	// ExecuteError processing
	var err error
	if be.DevError != model.DevErrorSuccess {
		query := be.GetDeviceReply(cmd)
		if be.Callback != nil {
			err = be.Callback.ExecuteError(ctx, &query)
		}
		if be.Log != nil {
			be.Log.Debug("Callback ExecuteError", slog.Any("query", query))
		}
	}
	return err
}

func (be *BaseEngine) RunStateChanged(ctx context.Context, state model.DevState) error {
	// StateChanged processing
	var err error
	if be.DevState != state {
		query := model.DeviceState{
			DeviceNotify: be.GetDeviceNotify(),
			StateNotify: model.StateNotify{
				OldState: be.DevState,
				NewState: state,
			},
		}
		if be.Callback != nil {
			err = be.Callback.StateChanged(ctx, &query)
		}
		if be.Log != nil {
			be.Log.Debug("Callback StateChanged", slog.Any("query", query))
		}
		be.DevState = state
	}
	return err
}

func (be *BaseEngine) RunActionPrompt(ctx context.Context, prompt model.DevPrompt) error {
	be.DevPrompt = prompt
	// ActionPrompt processing
	var err error
	if be.DevPrompt != model.DevPromptNone {
		query := &model.DevicePrompt{
			DeviceNotify: be.GetDeviceNotify(),
			PromptNotify: model.PromptNotify{
				Prompt: be.DevPrompt,
			},
		}
		if be.Callback != nil {
			err = be.Callback.ActionPrompt(ctx, query)
		}
		if be.Log != nil {
			be.Log.Debug("Callback ActionPrompt", slog.Any("query", query))
		}
	}
	return err
}

func (be *BaseEngine) RunReaderReturn(ctx context.Context, result string) error {
	be.DevResult = result
	// ReaderReturn processing
	var err error
	if be.DevResult != "" {
		query := &model.DeviceResult{
			DeviceNotify: be.GetDeviceNotify(),
			ResultNotify: model.ResultNotify{
				Result: be.DevResult,
			},
		}
		if be.Callback != nil {
			err = be.Callback.ReaderResult(ctx, query)
		}
		if be.Log != nil {
			be.Log.Debug("Callback ReaderResult", slog.Any("query", query))
		}
	}
	return err
}

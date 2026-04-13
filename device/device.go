package device

import (
	"context"
	"log/slog"

	"github.com/iftsoft/linker/model"
)

type Callback interface {
	model.DeviceCallback
}

type BaseEngine struct {
	Log       *slog.Logger
	DevName   string
	DevState  model.DevState
	DevError  model.DevError
	DevPrompt model.DevPrompt
	DevAction model.DevAction
	DevInform string
	DevReply  string
	Callback  Callback
}

func (be *BaseEngine) ClearDevice() {
	be.DevState = model.DevStateUndefined
	be.DevError = model.DevErrorSuccess
	be.DevPrompt = model.DevPromptNone
	be.DevAction = model.DevActionDoNothing
	be.DevInform = ""
	be.DevReply = ""
}

func (be *BaseEngine) RunDeviceReply(ctx context.Context, cmd string) error {
	// StateChanged processing
	var err error
	reply := &model.DeviceReply{
		Command: cmd,
		Device:  be.DevName,
		Action:  be.DevAction,
		State:   be.DevState,
		ErrCode: be.DevError,
		ErrText: be.DevReply,
	}
	if be.Callback != nil {
		err = be.Callback.DeviceReply(ctx, reply)
	}
	if be.Log != nil {
		be.Log.Debug("Callback DeviceReply: %s", reply.String())
	}
	return err
}

func (be *BaseEngine) RunExecuteError(ctx context.Context, errCode model.DevError, reason string) error {
	be.DevError = errCode
	be.DevReply = model.NewError(errCode, reason).Error()
	// ExecuteError processing
	var err error
	if be.DevError != model.DevErrorSuccess {
		query := &model.DeviceReply{
			Device:  be.DevName,
			Action:  be.DevAction,
			State:   be.DevState,
			ErrCode: be.DevError,
			ErrText: be.DevReply,
		}
		if be.Callback != nil {
			err = be.Callback.ExecuteError(ctx, query)
		}
		if be.Log != nil {
			be.Log.Debug("Callback ExecuteError: %s", query.String())
		}
	}
	return err
}

func (be *BaseEngine) RunStateChanged(ctx context.Context, state model.DevState) error {
	// StateChanged processing
	var err error
	if be.DevState != state {
		query := &model.DeviceState{
			Device:   be.DevName,
			Action:   be.DevAction,
			OldState: be.DevState,
			NewState: state,
		}
		if be.Callback != nil {
			err = be.Callback.StateChanged(ctx, query)
		}
		if be.Log != nil {
			be.Log.Debug("Callback StateChanged: %s", query.String())
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
			Device: be.DevName,
			Action: be.DevAction,
			Prompt: be.DevPrompt,
		}
		if be.Callback != nil {
			err = be.Callback.ActionPrompt(ctx, query)
		}
		if be.Log != nil {
			be.Log.Debug("Callback ActionPrompt: %s", query.String())
		}
	}
	return err
}

func (be *BaseEngine) RunReaderReturn(ctx context.Context, inform string) error {
	be.DevInform = inform
	// ReaderReturn processing
	var err error
	if be.DevInform != "" {
		query := &model.DeviceInform{
			Device: be.DevName,
			Action: be.DevAction,
			Inform: be.DevInform,
		}
		if be.Callback != nil {
			err = be.Callback.ReaderReturn(ctx, query)
		}
		if be.Log != nil {
			be.Log.Debug("Callback ReaderReturn: %s", query.String())
		}
	}
	return err
}

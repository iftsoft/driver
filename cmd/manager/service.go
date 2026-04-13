package main

import (
	"context"
	"log/slog"

	"github.com/iftsoft/linker/model"
)

type CallbackService struct {
	log *slog.Logger
}

func NewCallbackService(log *slog.Logger) *CallbackService {
	return &CallbackService{
		log: log,
	}
}

// SystemReply sends notification about device reply
func (cs *CallbackService) SystemReply(ctx context.Context, reply *model.SystemReply) error {
	cs.log.Info("CallbackService.SystemReply", "reply", reply.String())
	return nil
}

// SystemHealth sends notification about device reply
func (cs *CallbackService) SystemHealth(ctx context.Context, reply *model.SystemHealth) error {
	cs.log.Info("CallbackService.SystemHealth", "reply", reply.String())
	return nil
}

// DeviceReply sends notification about device reply
func (cs *CallbackService) DeviceReply(ctx context.Context, reply *model.DeviceReply) error {
	cs.log.Info("CallbackService.DeviceReply", "reply", reply.String())
	return nil
}

// ExecuteError sends notification about execute error
func (cs *CallbackService) ExecuteError(ctx context.Context, value *model.DeviceReply) error {
	cs.log.Info("CallbackService.ExecuteError", "value", value.String())
	return nil
}

// StateChanged sends notification about device state changing
func (cs *CallbackService) StateChanged(ctx context.Context, value *model.DeviceState) error {
	cs.log.Info("CallbackService.StateChanged", "value", value.String())
	return nil
}

// ActionPrompt sends notification about action prompt for user
func (cs *CallbackService) ActionPrompt(ctx context.Context, value *model.DevicePrompt) error {
	cs.log.Info("CallbackService.ActionPrompt", "value", value.String())
	return nil
}

// ReaderReturn sends notification about device reading result
func (cs *CallbackService) ReaderReturn(ctx context.Context, value *model.DeviceInform) error {
	cs.log.Info("CallbackService.ReaderReturn", "value", value.String())
	return nil
}

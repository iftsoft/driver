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

// GreetingInfo sends notification about device application
func (cs *CallbackService) GreetingInfo(ctx context.Context, value *model.GreetingInfo) error {
	cs.log.Info("CallbackService.GreetingInfo", slog.Any("value", value))
	return nil
}

// SystemReply sends notification about device reply
func (cs *CallbackService) SystemReply(ctx context.Context, value *model.SystemReply) error {
	cs.log.Info("CallbackService.SystemReply", slog.Any("value", value))
	return nil
}

// SystemHealth sends notification about device reply
func (cs *CallbackService) SystemHealth(ctx context.Context, value *model.SystemHealth) error {
	cs.log.Info("CallbackService.SystemHealth", slog.Any("value", value))
	return nil
}

// DeviceReply sends notification about device reply
func (cs *CallbackService) DeviceReply(ctx context.Context, value *model.DeviceReply) error {
	cs.log.Info("CallbackService.DeviceReply", slog.Any("value", value))
	return nil
}

// ExecuteError sends notification about execute error
func (cs *CallbackService) ExecuteError(ctx context.Context, value *model.DeviceReply) error {
	cs.log.Info("CallbackService.ExecuteError", slog.Any("value", value))
	return nil
}

// StateChanged sends notification about device state changing
func (cs *CallbackService) StateChanged(ctx context.Context, value *model.DeviceState) error {
	cs.log.Info("CallbackService.StateChanged", slog.Any("value", value))
	return nil
}

// ActionPrompt sends notification about action prompt for user
func (cs *CallbackService) ActionPrompt(ctx context.Context, value *model.DevicePrompt) error {
	cs.log.Info("CallbackService.ActionPrompt", slog.Any("value", value))
	return nil
}

// ReaderReturn sends notification about device reading result
func (cs *CallbackService) ReaderReturn(ctx context.Context, value *model.DeviceInform) error {
	cs.log.Info("CallbackService.ReaderReturn", slog.Any("value", value))
	return nil
}

// PrinterProgress sent notification about printing progress
func (cs *CallbackService) PrinterProgress(ctx context.Context, value *model.PrinterProgress) error {
	cs.log.Info("CallbackService.PrinterProgress", slog.Any("value", value))
	return nil
}

// CardPosition sends notification about new card position
func (cs *CallbackService) CardPosition(ctx context.Context, value *model.CardPosition) error {
	cs.log.Info("CallbackService.CardPosition", slog.Any("value", value))
	return nil
}

// CardDescription sends notification about card information
func (cs *CallbackService) CardDescription(ctx context.Context, value *model.CardDescription) error {
	cs.log.Info("CallbackService.CardDescription", slog.Any("value", value))
	return nil
}

// NoteAccepted sends notification about new note in escrow
func (cs *CallbackService) NoteAccepted(ctx context.Context, value *model.ValidatorAccept) error {
	cs.log.Info("CallbackService.NoteAccepted", slog.Any("value", value))
	return nil
}

// CashIsStored sends notification that note is stored to cassette
func (cs *CallbackService) CashIsStored(ctx context.Context, value *model.ValidatorAccept) error {
	cs.log.Info("CallbackService.CashIsStored", slog.Any("value", value))
	return nil
}

// CashReturned sends notification that note is returned to user
func (cs *CallbackService) CashReturned(ctx context.Context, value *model.ValidatorAccept) error {
	cs.log.Info("CallbackService.CashReturned", slog.Any("value", value))
	return nil
}

// ValidatorStore sends notification about current cassette state
func (cs *CallbackService) ValidatorStore(ctx context.Context, value *model.ValidatorBatch) error {
	cs.log.Info("CallbackService.ValidatorStore", slog.Any("value", value))
	return nil
}

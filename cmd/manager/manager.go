package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/iftsoft/linker/client/manager"
	"github.com/iftsoft/linker/model"
)

type Manager struct {
	log   *slog.Logger
	port  int
	name  string
	cli   *manager.ManagerClient
	setup model.DeviceSetup
}

func NewManager(log *slog.Logger, greeting *model.GreetingInfo) *Manager {
	return &Manager{
		log:  log,
		port: int(greeting.GrpcPort),
		name: greeting.DevName,
	}
}

func (cs *Manager) Process(ctx context.Context) {
	// manager client init
	address := fmt.Sprintf("127.0.0.1:%d", cs.port)
	managerCli, err := manager.NewManagerClient(ctx, cs.log, address)
	if err != nil {
		cs.log.Error("Manager client failed", slog.String("error", err.Error()))
		return
	}
	defer managerCli.Close()
	cs.cli = managerCli

	err = cs.RunDeviceFlow(ctx)
	if err != nil {
		cs.log.Error("Run device flow failed", slog.String("error", err.Error()))
		return
	}
}

func (cs *Manager) RunDeviceFlow(ctx context.Context) error {
	err := cs.RunSysStart(ctx)
	if err != nil {
		return fmt.Errorf("run sys start failed: %w", err)
	}

	err = cs.RunSysRestart(ctx)
	if err != nil {
		return fmt.Errorf("run sys restart failed: %w", err)
	}

	err = cs.RunSysHealth(ctx)
	if err != nil {
		return fmt.Errorf("run sys health failed: %w", err)
	}

	if cs.setup.Supported&model.ScopeFlagDevice == model.ScopeFlagDevice {
		err = cs.RunReset(ctx)
		if err != nil {
			return fmt.Errorf("run device reset failed: %w", err)
		}
		err = cs.RunStatus(ctx)
		if err != nil {
			return fmt.Errorf("run device status failed: %w", err)
		}
		err = cs.RunExecute(ctx)
		if err != nil {
			return fmt.Errorf("run execute operation failed: %w", err)
		}
		err = cs.RunCancel(ctx)
		if err != nil {
			return fmt.Errorf("run cancel execution failed: %w", err)
		}
	}

	if cs.setup.Supported&model.ScopeFlagPrinter == model.ScopeFlagPrinter {
		err = cs.RunInitPrinter(ctx)
		if err != nil {
			return fmt.Errorf("run reset printer failed: %w", err)
		}
		err = cs.RunPrintPage(ctx)
		if err != nil {
			return fmt.Errorf("run print page failed: %w", err)
		}
	}

	if cs.setup.Supported&model.ScopeFlagReader == model.ScopeFlagReader {
		err = cs.RunEnterCard(ctx)
		if err != nil {
			return fmt.Errorf("run enter failed: %w", err)
		}
		err = cs.RunReadCard(ctx)
		if err != nil {
			return fmt.Errorf("run read card failed: %w", err)
		}
		err = cs.RunEjectCard(ctx)
		if err != nil {
			return fmt.Errorf("run eject card failed: %w", err)
		}
		err = cs.RunCaptureCard(ctx)
		if err != nil {
			return fmt.Errorf("run capture card failed: %w", err)
		}
	}

	if cs.setup.Supported&model.ScopeFlagValidator == model.ScopeFlagValidator {
		err = cs.RunInitValidator(ctx)
		if err != nil {
			return fmt.Errorf("run init validator failed: %w", err)
		}
		err = cs.RunDoValidate(ctx)
		if err != nil {
			return fmt.Errorf("run do validate failed: %w", err)
		}
		err = cs.RunAcceptNote(ctx)
		if err != nil {
			return fmt.Errorf("run accept note failed: %w", err)
		}
		err = cs.RunReturnNote(ctx)
		if err != nil {
			return fmt.Errorf("run return note failed: %w", err)
		}
		err = cs.RunStopValidate(ctx)
		if err != nil {
			return fmt.Errorf("run stop validate failed: %w", err)
		}
	}

	err = cs.RunSysStop(ctx)
	if err != nil {
		return fmt.Errorf("run sys stop failed: %w", err)
	}

	err = cs.RunTerminate(ctx)
	if err != nil {
		return fmt.Errorf("run terminate failed: %w", err)
	}

	return nil
}

func (cs *Manager) RunSysStart(ctx context.Context) error {
	query := &model.ConfigUpdate{
		Device: cs.name,
	}
	reply, err := cs.cli.SysStart(ctx, query)
	if err != nil {
		cs.log.Error("Start device failed", slog.String("error", err.Error()))
	} else {
		cs.log.Debug("Start device done", slog.Any("query", query), slog.Any("reply", reply))
		cs.setup = reply.DeviceSetup
	}
	return err
}

func (cs *Manager) RunSysRestart(ctx context.Context) error {
	query := &model.ConfigUpdate{
		Device: cs.name,
	}
	reply, err := cs.cli.SysRestart(ctx, query)
	if err != nil {
		cs.log.Error("Restart device failed", slog.String("error", err.Error()))
	} else {
		cs.log.Debug("Restart device done", slog.Any("query", query), slog.Any("reply", reply))
		cs.setup = reply.DeviceSetup
	}
	return err
}

func (cs *Manager) RunSysHealth(ctx context.Context) error {
	query := &model.SystemQuery{
		Device: cs.name,
	}
	reply, err := cs.cli.SysHealth(ctx, query)
	if err != nil {
		cs.log.Error("Get device health failed", slog.String("error", err.Error()))
	} else {
		cs.log.Debug("Get device health done", slog.Any("query", query), slog.Any("reply", reply))
	}
	return err
}

func (cs *Manager) RunSysStop(ctx context.Context) error {
	query := &model.SystemQuery{
		Device: cs.name,
	}
	reply, err := cs.cli.SysStop(ctx, query)
	if err != nil {
		cs.log.Error("Stop device failed", slog.String("error", err.Error()))
	} else {
		cs.log.Debug("Stop device done", slog.Any("query", query), slog.Any("reply", reply))
	}
	return err
}

func (cs *Manager) RunTerminate(ctx context.Context) error {
	query := &model.SystemQuery{
		Device: cs.name,
	}
	reply, err := cs.cli.Terminate(ctx, query)
	if err != nil {
		cs.log.Error("Terminate app failed", slog.String("error", err.Error()))
	} else {
		cs.log.Debug("Terminate app done", slog.Any("query", query), slog.Any("reply", reply))
	}
	return err
}

func (cs *Manager) RunCancel(ctx context.Context) error {
	query := &model.DeviceQuery{
		Device: cs.name,
	}
	reply, err := cs.cli.Cancel(ctx, query)
	if err != nil {
		cs.log.Error("Cancel operation failed", slog.String("error", err.Error()))
	} else {
		cs.log.Debug("Cancel operation done", slog.Any("query", query), slog.Any("reply", reply))
	}
	return err
}

func (cs *Manager) RunReset(ctx context.Context) error {
	query := &model.DeviceQuery{
		Device: cs.name,
	}
	reply, err := cs.cli.Reset(ctx, query)
	if err != nil {
		cs.log.Error("Reset device failed", slog.String("error", err.Error()))
	} else {
		cs.log.Debug("Reset device done", slog.Any("query", query), slog.Any("reply", reply))
	}
	return err
}

func (cs *Manager) RunStatus(ctx context.Context) error {
	query := &model.DeviceQuery{
		Device: cs.name,
	}
	reply, err := cs.cli.Status(ctx, query)
	if err != nil {
		cs.log.Error("Get device status failed", slog.String("error", err.Error()))
	} else {
		cs.log.Debug("Get device status done", slog.Any("query", query), slog.Any("reply", reply))
	}
	return err
}

func (cs *Manager) RunExecute(ctx context.Context) error {
	query := &model.DeviceQuery{
		Device: cs.name,
	}
	reply, err := cs.cli.Execute(ctx, query)
	if err != nil {
		cs.log.Error("Execute device command failed", slog.String("error", err.Error()))
	} else {
		cs.log.Debug("Execute device command done", slog.Any("query", query), slog.Any("reply", reply))
	}
	return err
}

func (cs *Manager) RunInitPrinter(ctx context.Context) error {
	query := &model.PrinterSetup{
		Device: cs.name,
	}
	reply, err := cs.cli.InitPrinter(ctx, query)
	if err != nil {
		cs.log.Error("Init printer failed", slog.String("error", err.Error()))
	} else {
		cs.log.Debug("Init printer done", slog.Any("query", query), slog.Any("reply", reply))
	}
	return err
}

func (cs *Manager) RunPrintPage(ctx context.Context) error {
	query := &model.PrinterQuery{
		Device: cs.name,
	}
	reply, err := cs.cli.PrintPage(ctx, query)
	if err != nil {
		cs.log.Error("Print page failed", slog.String("error", err.Error()))
	} else {
		cs.log.Debug("Print page done", slog.Any("query", query), slog.Any("reply", reply))
	}
	return err
}

func (cs *Manager) RunEnterCard(ctx context.Context) error {
	query := &model.DeviceQuery{
		Device: cs.name,
	}
	reply, err := cs.cli.EnterCard(ctx, query)
	if err != nil {
		cs.log.Error("Enter card failed", slog.String("error", err.Error()))
	} else {
		cs.log.Debug("Enter card done", slog.Any("query", query), slog.Any("reply", reply))
	}
	return err
}

func (cs *Manager) RunEjectCard(ctx context.Context) error {
	query := &model.DeviceQuery{
		Device: cs.name,
	}
	reply, err := cs.cli.EjectCard(ctx, query)
	if err != nil {
		cs.log.Error("Eject card failed", slog.String("error", err.Error()))
	} else {
		cs.log.Debug("Eject card done", slog.Any("query", query), slog.Any("reply", reply))
	}
	return err
}

func (cs *Manager) RunCaptureCard(ctx context.Context) error {
	query := &model.DeviceQuery{
		Device: cs.name,
	}
	reply, err := cs.cli.CaptureCard(ctx, query)
	if err != nil {
		cs.log.Error("Capture card failed", slog.String("error", err.Error()))
	} else {
		cs.log.Debug("Capture card done", slog.Any("query", query), slog.Any("reply", reply))
	}
	return err
}

func (cs *Manager) RunReadCard(ctx context.Context) error {
	query := &model.DeviceQuery{
		Device: cs.name,
	}
	reply, err := cs.cli.ReadCard(ctx, query)
	if err != nil {
		cs.log.Error("Read card failed", slog.String("error", err.Error()))
	} else {
		cs.log.Debug("Read card done", slog.Any("query", query), slog.Any("reply", reply))
	}
	return err
}

func (cs *Manager) RunInitValidator(ctx context.Context) error {
	query := &model.ValidatorQuery{
		Device:   cs.name,
		Currency: model.CurrencyUSD,
	}
	reply, err := cs.cli.InitValidator(ctx, query)
	if err != nil {
		cs.log.Error("Init validator failed", slog.String("error", err.Error()))
	} else {
		cs.log.Debug("Init validator done", slog.Any("query", query), slog.Any("reply", reply))
	}
	return err
}

func (cs *Manager) RunDoValidate(ctx context.Context) error {
	query := &model.ValidatorQuery{
		Device:   cs.name,
		Currency: model.CurrencyUSD,
	}
	reply, err := cs.cli.DoValidate(ctx, query)
	if err != nil {
		cs.log.Error("Start validating failed", slog.String("error", err.Error()))
	} else {
		cs.log.Debug("Start validating done", slog.Any("query", query), slog.Any("reply", reply))
	}
	return err
}

func (cs *Manager) RunAcceptNote(ctx context.Context) error {
	query := &model.ValidatorQuery{
		Device:   cs.name,
		Currency: model.CurrencyUSD,
	}
	reply, err := cs.cli.AcceptNote(ctx, query)
	if err != nil {
		cs.log.Error("Accept note failed", slog.String("error", err.Error()))
	} else {
		cs.log.Debug("Accept note done", slog.Any("query", query), slog.Any("reply", reply))
	}
	return err
}

func (cs *Manager) RunReturnNote(ctx context.Context) error {
	query := &model.ValidatorQuery{
		Device:   cs.name,
		Currency: model.CurrencyUSD,
	}
	reply, err := cs.cli.ReturnNote(ctx, query)
	if err != nil {
		cs.log.Error("Return note failed", slog.String("error", err.Error()))
	} else {
		cs.log.Debug("Return note done", slog.Any("query", query), slog.Any("reply", reply))
	}
	return err
}

func (cs *Manager) RunStopValidate(ctx context.Context) error {
	query := &model.ValidatorQuery{
		Device:   cs.name,
		Currency: model.CurrencyUSD,
	}
	reply, err := cs.cli.StopValidate(ctx, query)
	if err != nil {
		cs.log.Error("Stop validating failed", slog.String("error", err.Error()))
	} else {
		cs.log.Debug("Stop validating done", slog.Any("query", query), slog.Any("reply", reply))
	}
	return err
}

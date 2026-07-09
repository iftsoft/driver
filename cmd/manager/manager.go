package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/iftsoft/linker/client/manager"
	"github.com/iftsoft/linker/model"
)

type Manager struct {
	log  *slog.Logger
	port int
	name string
	cli  *manager.ManagerClient
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
	sysCfg := &model.SystemConfig{
		Device: cs.name,
	}
	out, err := cs.cli.SysStart(ctx, sysCfg)
	if err != nil {
		cs.log.Error("Start device failed", slog.String("error", err.Error()))
	} else {
		cs.log.Debug("Start device done", slog.Any("reply", out))
	}

	sysQry := &model.SystemQuery{
		Device: cs.name,
	}
	out2, err := cs.cli.Terminate(ctx, sysQry)
	if err != nil {
		cs.log.Error("Terminate device failed", slog.String("error", err.Error()))
	} else {
		cs.log.Debug("Terminate device done", slog.Any("reply", out2))
	}

	return nil
}

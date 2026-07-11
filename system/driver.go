package system

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/iftsoft/linker/model"

	"github.com/iftsoft/driver/config"
	"github.com/iftsoft/driver/device"
)

type DeviceDriver struct {
	log       *slog.Logger
	config    config.DeviceConfig
	callback  device.Callback
	creator   device.DeviceCreator
	worker    device.DeviceWorker
	device    model.DeviceManager
	printer   model.PrinterManager
	reader    model.ReaderManager
	validator model.ValidatorManager
	settings  model.DeviceSetup
	mutex     sync.RWMutex
	done      chan struct{}
	devName   string
}

func NewDeviceDriver(setup *AppSetup, callback device.Callback, creator device.DeviceCreator) *DeviceDriver {
	dummy := NewNoopDevice(setup.Logger)
	drv := DeviceDriver{
		log:       setup.Logger,
		config:    setup.Config.Device,
		devName:   setup.Params.DevName,
		callback:  callback,
		creator:   creator,
		device:    dummy,
		printer:   dummy,
		reader:    dummy,
		validator: dummy,
		worker:    dummy,
	}
	return &drv
}

func (d *DeviceDriver) DeviceManager() model.DeviceManager {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.device
}

func (d *DeviceDriver) PrinterManager() model.PrinterManager {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.printer
}

func (d *DeviceDriver) ReaderManager() model.ReaderManager {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.reader
}

func (d *DeviceDriver) ValidatorManager() model.ValidatorManager {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.validator
}

func (d *DeviceDriver) CreateDevice(ctx context.Context, query *model.ConfigUpdate) error {
	if d.creator == nil {
		return errors.New("driver does not have a creator")
	}
	d.clearManagers()

	object, err := d.creator.CreateDevice(device.CreatorParams{
		DevName: d.devName, Logger: d.log, Config: &d.config, Callback: d.callback,
	})
	if err != nil {
		return fmt.Errorf("create device error: %w", err)
	}
	err = d.initManagers(object)
	if err != nil {
		return fmt.Errorf("init managers error: %w", err)
	}

	err = d.worker.StartDevice(ctx, query)
	if err != nil {
		return fmt.Errorf("start device error: %w", err)
	}

	d.done = make(chan struct{})
	go d.startDeviceLoop(context.Background())

	return nil
}

func (d *DeviceDriver) DeleteDevice(ctx context.Context) error {
	d.log.Info("Stopping system device")
	if d.done != nil {
		close(d.done)
		d.done = nil
	}
	defer d.clearManagers()

	err := d.worker.StopDevice(ctx)
	if err != nil {
		return fmt.Errorf("stop device error: %w", err)
	}

	return nil
}

func (d *DeviceDriver) CheckDevice(ctx context.Context) (*model.DeviceMetrics, error) {
	metrics, err := d.worker.CheckDevice(ctx)
	if err != nil {
		return nil, fmt.Errorf("check device error: %w", err)
	}
	return metrics, nil
}

func (d *DeviceDriver) clearManagers() {
	dummy := NewNoopDevice(d.log)

	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.device = dummy
	d.printer = dummy
	d.reader = dummy
	d.validator = dummy
	d.worker = dummy
}

func (d *DeviceDriver) initManagers(object interface{}) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.settings = model.DeviceSetup{}
	// Setup Worker driver interface
	if worker, ok := object.(device.DeviceWorker); ok {
		d.worker = worker
		d.settings = d.worker.DeviceSettings()
		d.settings.Supported = model.ScopeFlagSystem
	}
	// Setup Device scope interface
	if common, ok := object.(model.DeviceManager); ok {
		d.device = common
		d.settings.Supported |= model.ScopeFlagDevice
	}
	// Setup Printer scope interface
	if printer, ok := object.(model.PrinterManager); ok {
		d.printer = printer
		d.settings.Supported |= model.ScopeFlagPrinter
	}
	// Setup Reader scope interface
	if reader, ok := object.(model.ReaderManager); ok {
		d.reader = reader
		d.settings.Supported |= model.ScopeFlagReader
	}
	// Setup Validator scope interface
	if validator, ok := object.(model.ValidatorManager); ok {
		d.validator = validator
		d.settings.Supported |= model.ScopeFlagValidator
	}

	validMask := model.ScopeFlagSystem | model.ScopeFlagDevice
	if d.settings.Supported&validMask == validMask {
		return nil
	}
	return errors.New("device object is not valid")
}

func (d *DeviceDriver) startDeviceLoop(ctx context.Context) {
	d.log.Debug("System device loop is started", "device", d.devName)
	defer d.log.Debug("System device loop is stopped", "device", d.devName)

	tick := time.NewTicker(100 * time.Millisecond)
	defer tick.Stop()

	for {
		select {
		case <-d.done:
			return
		case tm := <-tick.C:
			d.log.Debug("System device %s onTimerTick", "device", d.devName, "moment", tm.Format(time.StampMilli))
			if d.worker != nil {
				_ = d.worker.DeviceTimer(ctx, tm.Unix())
			}
		}
	}
}

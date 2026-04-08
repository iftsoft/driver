package linker

import (
	"errors"
	"log/slog"

	"github.com/iftsoft/driver/config"
	//	"github.com/iftsoft/device/core"
)

var (
	errPortNotOpen = errors.New("port is not open")
)

const (
	linkerBufferSize = 1024
)

type PortReader interface {
	OnRead(data []byte) int
}

type PortLinker interface {
	Open() error
	Close() error
	Flash() error
	IsOpen() bool
	Write(data []byte) (int, error)
}

func GetPortLinker(log *slog.Logger, cfg *config.LinkerConfig, call PortReader) PortLinker {
	dummy := NewDummyLink(log, call)
	if cfg == nil {
		return dummy
	}
	switch cfg.LinkType {
	case config.LinkTypeNone:
		return dummy
	case config.LinkTypeSerial:
		return NewSerialLink(log, cfg.Serial, call)
	case config.LinkTypeHidUsb:
		return NewDummyLinker(log, cfg.HidUsb, call)
	}
	return dummy
}

func GetLinkerPorts(log *slog.Logger) error {
	_, err := EnumerateSerialPorts(log)
	if err != nil {
		log.Error("Serial port error: %s", err.Error())
	}
	_, err = EnumerateHidUsbPorts(log)
	if err != nil {
		log.Error("HidUsb port error: %s", err.Error())
	}
	return err
}

package hardware

import (
	"log/slog"
)

type DummyLink struct {
	log    *slog.Logger
	reader PortReader
	isOpen bool
}

func NewDummyLink(log *slog.Logger, call PortReader) *DummyLink {
	d := DummyLink{
		log:    log,
		reader: call,
		isOpen: false,
	}
	return &d
}

func (d *DummyLink) Open() error {
	d.isOpen = true
	d.log.Debug("DummyLink run cmd:Open")
	return nil
}

func (d *DummyLink) Close() error {
	d.isOpen = false
	d.log.Debug("DummyLink run cmd:Close")
	return nil
}

func (d *DummyLink) Flash() error {
	d.log.Debug("DummyLink run cmd:Flash")
	return nil
}

func (d *DummyLink) IsOpen() bool {
	return d.isOpen
}

func (d *DummyLink) Write(data []byte) (int, error) {
	if !d.isOpen {
		return 0, errPortNotOpen
	}
	d.log.Debug("DummyLink write data", slog.Any("data", data))
	go func(dump []byte) {
		if d.reader != nil {
			d.log.Debug("DummyLink read data", slog.Any("dump", dump))
			d.reader.OnRead(dump)
		}
	}(data)
	return len(data), nil
}

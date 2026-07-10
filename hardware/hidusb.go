package hardware

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/karalabe/hid"

	"github.com/iftsoft/driver/config"
)

type HidUsbLink struct {
	config *config.HidUsbConfig
	log    *slog.Logger
	link   *hid.Device
	reader PortReader
	isOpen bool
}

func NewDummyLinker(log *slog.Logger, cfg *config.HidUsbConfig, call PortReader) *HidUsbLink {
	h := HidUsbLink{
		config: cfg,
		log:    log,
		link:   nil,
		reader: call,
		isOpen: false,
	}
	return &h
}

func EnumerateHidUsbPorts(log *slog.Logger) (list []*config.HidUsbConfig, err error) {
	//	defer core.PanicRecover(&err, out)
	log.Debug("HidUsb port enumeration")
	units := hid.Enumerate(0, 0)
	if units == nil {
		return nil, errors.New("hidapi library is not working")
	}
	for i, unit := range units {
		line := fmt.Sprintf("   Port#%d - %d:%d/%s (%s - %s, %s)", i,
			unit.VendorID, unit.ProductID, unit.Serial, unit.Manufacturer, unit.Product, unit.Path)
		log.Debug(line)
		item := &config.HidUsbConfig{
			VendorID:  uint32(unit.VendorID),
			ProductID: uint32(unit.ProductID),
			SerialNo:  unit.Serial,
		}
		list = append(list, item)
	}
	return list, err
}

func (h *HidUsbLink) Open() (err error) {
	//	defer core.PanicRecover(&err, h.log)
	if h.config == nil {
		return errors.New("HidUsb config is not set")
	}
	info := hid.DeviceInfo{
		VendorID:  uint16(h.config.VendorID),
		ProductID: uint16(h.config.ProductID),
		Serial:    h.config.SerialNo,
	}
	h.link, err = info.Open()
	if err == nil {
		go h.readingLoop()
	}
	return err
}

func (h *HidUsbLink) Close() (err error) {
	//defer core.PanicRecover(&err, h.log)
	if h.link == nil {
		return err
	}
	err = h.link.Close()
	if err == nil {
		h.link = nil
	}
	return err
}

func (h *HidUsbLink) Flash() error {
	if h.link == nil {
		return errPortNotOpen
	}
	return nil
}

func (h *HidUsbLink) IsOpen() bool {
	return h.isOpen
}

func (h *HidUsbLink) Write(data []byte) (n int, err error) {
	//defer core.PanicRecover(&err, h.log)
	if h.link == nil {
		return 0, errPortNotOpen
	}
	n, err = h.link.Write(data)
	line := fmt.Sprintf("Write to hidapi port %d:%d return %v", h.config.VendorID, h.config.ProductID, err)
	h.log.Debug(line)
	return n, err
}

func (h *HidUsbLink) readData(data []byte) (n int, err error) {
	//defer core.PanicRecover(&err, h.log)
	if h.link == nil {
		return 0, errPortNotOpen
	}
	n, err = h.link.Read(data)
	line := fmt.Sprintf("Read from hidapi port %d:%d return %v",
		h.config.VendorID, h.config.ProductID, err)
	h.log.Debug(line)
	return n, err
}

func (h *HidUsbLink) readingLoop() {
	h.isOpen = true
	defer func() { h.isOpen = false }()
	h.log.Debug("HidUsb reading loop is started")
	defer h.log.Debug("HidUsb reading loop is stopped")

	rest := []byte{}
	buff := make([]byte, linkerBufferSize)
	for {
		n, err := h.readData(buff)
		if n > 0 {
			dump := buff[0:n]
			data := append(rest, dump...)
			rest = h.processData(data)
		}
		if err != nil {
			h.log.Warn("HidUsb ReadData error", slog.String("error", err.Error()))
			return
		}
	}
}

func (h *HidUsbLink) processData(data []byte) (out []byte) {
	if h.reader == nil {
		return nil
	}
	k := h.reader.OnRead(data)
	if k == 0 {
		return data
	}
	if k > 0 && k < len(data) {
		return h.processData(data[k:])
	}
	return nil
}

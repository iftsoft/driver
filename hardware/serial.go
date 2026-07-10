package hardware

import (
	"errors"
	"fmt"
	"log/slog"

	"go.bug.st/serial"

	"github.com/iftsoft/driver/config"
)

type SerialLink struct {
	config *config.SerialConfig
	log    *slog.Logger
	port   serial.Port
	reader PortReader
	isOpen bool
}

func NewSerialLink(log *slog.Logger, cfg *config.SerialConfig, call PortReader) *SerialLink {
	s := SerialLink{
		config: cfg,
		log:    log,
		port:   nil,
		reader: call,
		isOpen: false,
	}
	return &s
}

func EnumerateSerialPorts(log *slog.Logger) (list []string, err error) {
	//defer core.PanicRecover(&err, log)
	log.Debug("Serial port enumeration")
	list, err = serial.GetPortsList()
	for i, ser := range list {
		line := fmt.Sprintf("   Port#%d - %s", i, ser)
		log.Debug(line)
	}
	return list, err
}

func (s *SerialLink) Open() (err error) {
	//defer core.PanicRecover(&err, s.log)
	if s.config == nil {
		return errors.New("serial config is not set")
	}
	m := &serial.Mode{
		BaudRate: int(s.config.BaudRate),
		DataBits: int(s.config.DataBits),
		Parity:   serial.Parity(s.config.Parity),
		StopBits: serial.StopBits(s.config.StopBits),
	}
	s.port, err = serial.Open(s.config.PortName, m)
	if s.port == nil && err == nil {
		err = errPortNotOpen
	}
	if err == nil {
		go s.readingLoop()
	} else {
		s.port = nil
	}
	//s.log.Debug("Open serial port", s.config.PortName, core.GetErrorText(err))
	return err
}

func (s *SerialLink) Close() (err error) {
	//defer core.PanicRecover(&err, s.log)
	if s.port == nil {
		return err
	}
	err = s.port.Close()
	//s.log.Debug("Close serial port %s return %s", s.config.PortName, core.GetErrorText(err))
	if err == nil {
		s.port = nil
	}
	return err
}

func (s *SerialLink) Flash() (err error) {
	//defer core.PanicRecover(&err, s.log)
	if s.port == nil {
		err = errPortNotOpen
	}
	if err == nil {
		err = s.port.ResetInputBuffer()
	}
	if err == nil {
		err = s.port.ResetOutputBuffer()
	}
	//s.log.Debug("Flash serial port %s return %s", s.config.PortName, core.GetErrorText(err))
	return err
}

func (s *SerialLink) IsOpen() bool {
	return s.isOpen
}

func (s *SerialLink) Write(data []byte) (n int, err error) {
	//defer core.PanicRecover(&err, s.log)
	if s.port == nil {
		return 0, errPortNotOpen
	}
	//s.log.Debug("Serial write data : %s", core.GetBinaryDump(data))
	n, err = s.port.Write(data)
	//s.log.Debug("Write to serial port %s return %s", s.config.PortName, core.GetErrorText(err))
	return n, err
}

func (s *SerialLink) readData(data []byte) (n int, err error) {
	//defer core.PanicRecover(&err, s.log)
	if s.port == nil {
		return 0, errPortNotOpen
	}
	n, err = s.port.Read(data)
	//s.log.Debug("Serial read data : %s", core.GetBinaryDump(data[0:n]))
	//s.log.Debug("Read from serial port %s of %d bytes return %s",
	//	s.config.PortName, n, core.GetErrorText(err))
	return n, err
}

func (s *SerialLink) readingLoop() {
	s.isOpen = true
	defer func() { s.isOpen = false }()
	s.log.Debug("Serial reading loop is started")
	defer s.log.Debug("Serial reading loop is stopped")

	rest := []byte{}
	for {
		buff := make([]byte, linkerBufferSize)
		n, err := s.readData(buff)
		if n > 0 {
			dump := buff[0:n]
			data := append(rest, dump...)
			rest = s.processData(data)
		}
		if err != nil {
			s.log.Warn("Serial ReadData failed", slog.String("error", err.Error()))
			return
		}
	}
}

func (s *SerialLink) processData(data []byte) (out []byte) {
	//s.log.Debug("Process reply data : %s", core.GetBinaryDump(data))
	if s.reader == nil {
		return nil
	}
	if len(data) == 0 {
		return nil
	}
	k := s.reader.OnRead(data)
	if k == 0 {
		return data
	}
	if k > 0 && k < len(data) {
		return s.processData(data[k:])
	}
	return nil
}

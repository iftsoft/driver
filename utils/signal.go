package utils

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var signalChan chan os.Signal

func init() {
	signalChan = make(chan os.Signal)
}

func WaitForSignal(out *slog.Logger) {
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	s := <-signalChan
	if out != nil {
		out.Info("Got signal: %v, exiting.", s)
	}
}

func SendQuitSignal(wait int) {
	go func() {
		time.Sleep(time.Millisecond * time.Duration(wait))
		signalChan <- syscall.SIGTERM
	}()
}

package utils

/*
import (
	"log"
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
	time.Sleep(time.Millisecond * time.Duration(wait))

	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		log.Fatal(err)
	}

	// On a Unix-like system, pressing Ctrl+C on a keyboard sends a
	// SIGINT signal to the process of the program in execution.
	//
	// This example simulates that by sending a SIGINT signal to itself.
	if err = p.Signal(os.Interrupt); err != nil {
		log.Fatal(err)
	}
}
*/

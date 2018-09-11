package main

import (
	"os"
	"os/signal"
	"syscall"

	_ "github.com/air-examples/chatroom/handlers"
	"github.com/sheng/air"
)

func main() {
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		if err := air.Serve(); err != nil {
			air.ERROR(err.Error())
		}
	}()

	<-shutdownChan
	air.INFO("shutting down the server")
	air.Shutdown(0)
	air.INFO("server gracefully stopped")
}

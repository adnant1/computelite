package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

// waitForShutdown blocks until a termination signal is received
func waitForShutdown() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh
	log.Println("[app] shutdown signal received")
}

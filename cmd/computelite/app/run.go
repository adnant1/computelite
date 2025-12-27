package app

import (
	"log"

	"github.com/adnant1/computelite/pkg/cluster"
)

// Run initializes and starts the computelite simulation
func Run() int {
	log.Println("[app] starting computelite")

	// Core cluster state
	cs := cluster.NewCluster()

	// Stop channel (shared lifecycle)
	stop := make(chan struct{})

	// Start controllers
	startControllers(cs, stop)

	// Start continuously printing cluster state
	go startSnapshotPrinter(cs, stop)

	// Run bootstrap scenario
	go runScenario(cs, stop)

	// Wait for shutdown
	waitForShutdown()
	close(stop)

	log.Println("[app] computelite stopped cleanly")
	return 0
}

package app

import (
	"time"

	"github.com/adnant1/computelite/pkg/cluster"
)

// startSnapshotPrinter periodically prints the cluster state snapshot
func startSnapshotPrinter(cs *cluster.ClusterState, stop <-chan struct{}) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cs.PrintSnapshot()
		case <-stop:
			return
		}
	}
}

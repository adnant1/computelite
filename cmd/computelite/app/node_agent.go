package app

import (
	"log"
	"time"

	"github.com/adnant1/computelite/pkg/cluster"
)

// startNodeAgent simulates a node-level agent (like kubelet)
// that periodically reports heartbeats to the control plane.
func startNodeAgent(
	cs *cluster.ClusterState,
	nodeID string,
	heartbeatInterval time.Duration,
	stop <-chan struct{},
) {
	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()

	log.Printf("[node-agent] node=%s started (heartbeat=%s)", nodeID, heartbeatInterval)

	for {
		select {
		case <-ticker.C:
			if err := cs.RecordHeartbeat(nodeID); err != nil {
				log.Printf(
					"[node-agent] node=%s failed to record heartbeat: %v",
					nodeID,
					err,
				)
				return
			}

		case <-stop:
			log.Printf("[node-agent] node=%s stopping", nodeID)
			return
		}
	}
}

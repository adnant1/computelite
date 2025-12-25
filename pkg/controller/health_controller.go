package controller

import (
	"time"

	"github.com/adnant1/computelite/pkg/cluster"
)

// HealthController is responsible for monitoring and updating
// the health status of nodes in the cluster.
type HealthController struct {
	clusterState      *cluster.ClusterState
	reconcileInterval time.Duration
}

// NewHealthController configures a health controller with the given cluster state and reconciliation interval.
func NewHealthController(cs *cluster.ClusterState, interval time.Duration) *HealthController {
	return &HealthController{
		clusterState:      cs,
		reconcileInterval: interval,
	}
}

// Run begins monitoring and updating
func (hc *HealthController) Run() {
	ticker := time.NewTicker(hc.reconcileInterval)
	defer ticker.Stop()

	// will soon add stop channel handling
	for {
		select {
		case <-ticker.C:
			hc.reconcile()
		}
	}
}

// reconcile checks the health of each node and updates their status accordingly
func (hc *HealthController) reconcile() {
	now := time.Now()

	for nodeID, node := range hc.clusterState.Nodes {
		desiredHealth := cluster.EvaluateNodeHealth(node.LastHeartbeat, now)

		if node.Health != desiredHealth {
			hc.clusterState.UpdateNodeHealth(nodeID, desiredHealth)
		}
	}
}
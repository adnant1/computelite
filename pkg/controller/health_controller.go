package controller

import (
	"log"
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
func (hc *HealthController) Run(stop <-chan struct{}) {
	ticker := time.NewTicker(hc.reconcileInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			hc.reconcile()
		case <-stop:
			log.Printf("[health-controller] stopping health controller")
			return
		}
	}
}

// reconcile checks the health of each node and updates their status accordingly
func (hc *HealthController) reconcile() {
	now := time.Now()
	nodes := hc.clusterState.ListNodes()

	for _, node := range nodes {
		desiredHealth := cluster.EvaluateNodeHealth(node.LastHeartbeat, now)

		if node.Health != desiredHealth {
			hc.clusterState.UpdateNodeHealth(node.ID, desiredHealth)
			log.Printf("[health-controller] node=%s health changed to %s", node.ID, desiredHealth)
		}
	}
}
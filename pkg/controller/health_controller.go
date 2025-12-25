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

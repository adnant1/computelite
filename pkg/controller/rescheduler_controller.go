package controller

import (
	"time"

	"github.com/adnant1/computelite/pkg/cluster"
)

// ReschedulerController watches and manages the rescheduling of jobs on unhealthy nodes
type ReschedulerController struct {
	clusterState      *cluster.ClusterState
	reconcileInterval time.Duration
}

// NewReschedulerController configures a rescheduler controller with the given cluster state and reconciliation interval
func NewReschedulerController(cs *cluster.ClusterState, interval time.Duration) *ReschedulerController {
	return &ReschedulerController{
		clusterState:      cs,
		reconcileInterval: interval,
	}
}

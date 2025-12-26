package controller

import (
	"time"

	"github.com/adnant1/computelite/pkg/cluster"
	"github.com/adnant1/computelite/pkg/scheduler"
)

// SchedulerController watches for unscheduled jobs and assigns them to nodes
type SchedulerController struct {
	clusterState      *cluster.ClusterState
	policy            scheduler.SchedulerPolicy
	reconcileInterval time.Duration
}

// NewSchedulerController configures a scheduler controller with the given cluster state, scheduling policy, and reconciliation interval
func NewSchedulerController(clusterState *cluster.ClusterState, policy scheduler.SchedulerPolicy, reconcileInterval time.Duration) *SchedulerController {
	return &SchedulerController{
		clusterState:      clusterState,
		policy:            policy,
		reconcileInterval: reconcileInterval,
	}
}
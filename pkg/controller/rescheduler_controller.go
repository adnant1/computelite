package controller

import (
	"log"
	"time"

	"github.com/adnant1/computelite/pkg/api"
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

// Run begins watching and managing job rescheduling
func (rc *ReschedulerController) Run(stop <-chan struct{}) {
	ticker := time.NewTicker(rc.reconcileInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rc.reconcile()
		case <-stop:
			log.Printf("[rescheduler-controller] stopping rescheduler controller")
			return
		}
	}
}

// reconcile checks for jobs on unhealthy nodes and reschedules them
func (rc *ReschedulerController) reconcile() {
	for jobID, job := range rc.clusterState.Jobs {
		if job.AssignedNodeID != "" {
			node, exists := rc.clusterState.Nodes[job.AssignedNodeID]

			if exists && node.Health == api.Unhealthy {
				rc.clusterState.EvictAndRequeueJob(jobID)
				log.Printf("[rescheduler-controller] job=%d evicted from node=%s (node unhealthy)", jobID, job.AssignedNodeID)
			}
		}
	}
}
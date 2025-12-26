package controller

import (
	"log"
	"time"

	"github.com/adnant1/computelite/pkg/api"
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

// Run begins watching and scheduling jobs
func (sc *SchedulerController) Run() {
	ticker := time.NewTicker(sc.reconcileInterval)
	defer ticker.Stop()

	// will soon add stop channel handling
	for {
		select {
		case <-ticker.C:
			sc.reconcile()
		}
	}
}

// reconcile checks for a single pending job and attempts to schedule it onto a healthy node
func (sc *SchedulerController) reconcile() {
	healthyNodes := make(map[string]*api.Node)
	for nodeID, node := range sc.clusterState.Nodes {
		if node.Health == api.Healthy {
			healthyNodes[nodeID] = node
		}
	}

	if len(healthyNodes) == 0 {
		log.Printf("[scheduler-controller] no healthy nodes available for scheduling")
		return
	}

	for jobID, job := range sc.clusterState.Jobs {
		if job.State == api.Pending {
			nodeID, ok := sc.policy.SelectNode(job, healthyNodes)

			if !ok {
				log.Printf("[scheduler-controller] job=%d pending: no suitable node found", jobID)
				return
			}

			err := sc.clusterState.AssignJob(jobID, nodeID)
			if err != nil {
				log.Printf("[scheduler-controller] job=%d failed to assign to node=%s: %v", jobID, nodeID, err)
				return
			}

			log.Printf("[scheduler-controller] job=%d assigned to node=%s (cpu=%d, mem=%d)", 
			jobID, 
			nodeID, 
			job.Requires.CPU, 
			job.Requires.Memory)
			return
		}
	}
}
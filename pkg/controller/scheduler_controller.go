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
func (sc *SchedulerController) Run(stop <-chan struct{}) {
	ticker := time.NewTicker(sc.reconcileInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			sc.reconcile()
		case <-stop:
			log.Printf("[scheduler-controller] stopping scheduler controller")
			return
		}
	}
}

// reconcile checks for a single pending job and attempts to schedule it onto a healthy node
func (sc *SchedulerController) reconcile() {
	healthyNodes := sc.clusterState.ListHealthyNodes()
	pendingJobs := sc.clusterState.ListJobsByState(api.Pending)

	if len(healthyNodes) == 0 {
		log.Printf("[scheduler-controller] no healthy nodes available for scheduling")
		return
	}

	for _, job := range pendingJobs {
		nodeID, ok := sc.policy.SelectNode(job, healthyNodes)

		if !ok {
			log.Printf("[scheduler-controller] job=%d pending: no suitable node found", job.ID)
			return
		}

		err := sc.clusterState.AssignJob(job.ID, nodeID)
		if err != nil {
			log.Printf("[scheduler-controller] job=%d failed to assign to node=%s: %v", job.ID, nodeID, err)
			return
		}

		log.Printf("[scheduler-controller] job=%d assigned to node=%s (cpu=%d, mem=%d)", 
		job.ID, 
		nodeID, 
		job.Requires.CPU, 
		job.Requires.Memory)
		return
		
	}
}
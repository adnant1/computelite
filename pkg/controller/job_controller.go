package controller

import (
	"log"
	"time"

	"github.com/adnant1/computelite/pkg/api"
	"github.com/adnant1/computelite/pkg/cluster"
)

// JobController watches and manages the lifecycle of jobs in the cluster
type JobController struct {
	clusterState      *cluster.ClusterState
	reconcileInterval time.Duration
}

// NewJobController configures a job controller with the given cluster state and reconciliation interval
func NewJobController(cs *cluster.ClusterState, interval time.Duration) *JobController {
	return &JobController{
		clusterState:      cs,
		reconcileInterval: interval,
	}
}

// Run begins watching and managing jobs
func (jc *JobController) Run(stop <-chan struct{}) {
	ticker := time.NewTicker(jc.reconcileInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			jc.reconcile()
		case <-stop:
			log.Printf("[job-controller] stopping job controller")
			return
		}
	}
}

// reconcile checks the state of each job and updates their status accordingly
func (jc *JobController) reconcile() {
	now := time.Now()

	for jobID, job := range jc.clusterState.Jobs {
		switch job.State {
		case api.Assigned:
			if job.AssignedNodeID != "" {
				jc.clusterState.UpdateJobState(jobID, api.Running)
				log.Printf("[job-controller] job=%d assigned -> running (node=%s)", jobID, job.AssignedNodeID)
			}
		case api.Running:
			if now.Sub(job.StartedAt) >= api.RunDuration {
				jc.clusterState.UpdateJobState(jobID, api.Succeeded)
				log.Printf("[job-controller] job=%d running -> succeeded (duration=%s)", jobID, api.RunDuration)
			}
		// other states do not require action
		default:
			continue
		}
	}
}


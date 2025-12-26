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

// reconcile updates job states accordingly
func (jc *JobController) reconcile() {
	now := time.Now()
	assignedJobs := jc.clusterState.ListJobsByState(api.Assigned)
	runningJobs := jc.clusterState.ListJobsByState(api.Running)

	for _, job := range assignedJobs {
		if job.AssignedNodeID != "" {
			jc.clusterState.UpdateJobState(job.ID, api.Running)
			log.Printf("[job-controller] job=%d assigned -> running (node=%s)", job.ID, job.AssignedNodeID)
		}
	}

	for _, job := range runningJobs {
		if now.Sub(job.StartedAt) >= api.RunDuration {
			jc.clusterState.UpdateJobState(job.ID, api.Succeeded)
			log.Printf("[job-controller] job=%d running -> succeeded (duration=%s)", job.ID, api.RunDuration)
		}
	}
}


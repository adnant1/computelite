package controller

import (
	"time"

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


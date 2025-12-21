package scheduler

import "github.com/adnant1/computelite/pkg/cluster"

type Scheduler struct {
	Cluster     *cluster.ClusterState // Reference to the cluster the scheduler manages
	PendingJobs *JobQueue             // List of pending jobs waiting to be scheduled
}
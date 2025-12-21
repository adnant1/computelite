package scheduler

import "github.com/adnant1/computelite/pkg/cluster"

type Scheduler struct {
	Cluster     *cluster.ClusterState // Reference to the cluster the scheduler manages
	PendingJobs *JobQueue             // List of pending jobs waiting to be scheduled
}

// SubmitJob adds a new job to the pending jobs queue
// does not schedule any jobs, only adds to the JobQueue
func (s *Scheduler) SubmitJob(job *Job) {
	job.State = Pending
	s.PendingJobs.Enqueue(job)
}
package scheduler

import "github.com/adnant1/computelite/pkg/cluster"

type Scheduler struct {
	Cluster     *cluster.ClusterState // Reference to the cluster the scheduler manages
	PendingJobs *JobQueue             // List of pending jobs waiting to be scheduled
}

// NewScheduler initializes and returns a new Scheduler
func NewScheduler(cluster *cluster.ClusterState) *Scheduler {
	return &Scheduler{
		Cluster:     cluster,
		PendingJobs: &JobQueue{},
	}
}

// SubmitJob adds a new job to the pending jobs queue
// does not schedule any jobs, only adds to the JobQueue
func (s *Scheduler) SubmitJob(job *Job) {
	job.State = Pending
	s.PendingJobs.Enqueue(job)
}

// ScheduleOne attempts to schedule one pending job onto a suitable node
// for simplicity, it uses a first-fit algorithm
func (s *Scheduler) ScheduleOne() bool {
	job, ok := s.PendingJobs.Dequeue()
	if !ok {
		return false // no pending jobs to schedule
	}

	for _, node := range s.Cluster.Nodes {
		availableCPU := node.TotalCapacity.CPU - node.Allocated.CPU
		availableMemory := node.TotalCapacity.Memory - node.Allocated.Memory

		if availableCPU >= job.Requires.CPU && availableMemory >= job.Requires.Memory {
			// allocate resources on the node
			node.Allocated.CPU += job.Requires.CPU
			node.Allocated.Memory += job.Requires.Memory

			s.Cluster.RunningJobs[job.ID] = node
			job.State = Running

			return true
		}
	}
	s.PendingJobs.Enqueue(job)
	return false // no suitable node found
}

// ScheduleAll attempts to schedule all pending jobs
// only runs a single pass through the pending jobs
func (s *Scheduler) ScheduleAll() {
	for {
		if !s.ScheduleOne() {
			break
		}
	}
}
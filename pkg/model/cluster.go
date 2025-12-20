package model

type ClusterState struct {
	Nodes		map[string]*Node // Mapping of node IDs to Node structs
	PendingJobs	*JobQueue        // List of pending jobs in the cluster
	RunningJobs	map[int64]*Node  // Mapping of running job IDs to their assigned nodes
}

// SubmitJob adds a new job to the pending jobs queue
// for now assume all invariants are true
// does not schedule any jobs, only adds to the JobQueue
func (cs *ClusterState) SubmitJob(job *Job) {
	cs.PendingJobs.Enqueue(job)
}
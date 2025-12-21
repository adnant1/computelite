package cluster

import "github.com/adnant1/computelite/pkg/scheduler"

type ClusterState struct {
	Nodes		map[string]*Node // Mapping of node IDs to Node structs
	PendingJobs	*JobQueue        // List of pending jobs in the cluster
	RunningJobs	map[int64]*Node  // Mapping of running job IDs to their assigned nodes
}

// SubmitJob adds a new job to the pending jobs queue
// for now assume all invariants are true
// does not schedule any jobs, only adds to the JobQueue
func (cs *ClusterState) SubmitJob(job *scheduler.Job) {
	cs.PendingJobs.Enqueue(job)
}

// AddNode adds a new node to the cluster state
func (cs *ClusterState) AddNode(node *Node) {
	cs.Nodes[node.ID] = node
}
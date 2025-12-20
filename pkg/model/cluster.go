package model

type ClusterState struct {
	Nodes		map[string]*Node // Mapping of node IDs to Node structs
	PendingJobs	[]*Job           // List of pending jobs in the cluster
	RunningJobs	map[int64]*Node  // Mapping of running job IDs to their assigned nodes
}
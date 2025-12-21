package cluster

type ClusterState struct {
	Nodes		map[string]*Node // Mapping of node IDs to Node structs
	RunningJobs	map[int64]*Node  // Mapping of running job IDs to their assigned nodes
}

// AddNode adds a new node to the cluster state
func (cs *ClusterState) AddNode(node *Node) {
	cs.Nodes[node.ID] = node
}
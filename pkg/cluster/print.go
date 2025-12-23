package cluster

import "fmt"

// PrintSummary prints a high-level overview of the cluster state.
func (cs *ClusterState) PrintSummary() {
	fmt.Println("=== Cluster Summary ===")
	fmt.Printf("Nodes: %d\n", len(cs.Nodes))
	fmt.Printf("Running Jobs: %d\n", len(cs.RunningJobs))
}

// PrintNodeUtilization prints resource usage for each node.
func (cs *ClusterState) PrintNodeUtilization() {
	fmt.Println("\n=== Node Utilization ===")

	for _, node := range cs.Nodes {
		fmt.Printf("Node %s:\n", node.ID)
		fmt.Printf("  CPU:    %d / %d\n",
			node.Allocated.CPU,
			node.TotalCapacity.CPU,
		)
		fmt.Printf("  Memory: %d / %d\n",
			node.Allocated.Memory,
			node.TotalCapacity.Memory,
		)
	}
}

// PrintRunningJobs prints all running jobs and their assigned nodes.
func (cs *ClusterState) PrintRunningJobs() {
	fmt.Println("\n=== Running Jobs ===")

	if len(cs.RunningJobs) == 0 {
		fmt.Println("No running jobs")
		return
	}

	for jobID, node := range cs.RunningJobs {
		fmt.Printf("Job %d -> %s\n", jobID, node.ID)
	}
}

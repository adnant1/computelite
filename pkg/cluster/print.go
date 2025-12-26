package cluster

import (
	"fmt"

	"github.com/adnant1/computelite/pkg/api"
)

// PrintSnapshot prints a snapshot of the current cluster state
func (cs *ClusterState) PrintSnapshot() {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	fmt.Println("\n================ CLUSTER SNAPSHOT ================")
	cs.printJobSummaryLocked()
	cs.printNodeUtilizationLocked()
	cs.printNodeHealthLocked()
	fmt.Println("==================================================")
}

// printJobSummaryLocked prints a summary of jobs in the cluster
func (cs *ClusterState) printJobSummaryLocked() {
	counts := map[api.JobState]int{}

	for _, job := range cs.Jobs {
		counts[job.State]++
	}

	fmt.Println("\nJobs:")
	fmt.Printf(
		"  Pending=%d Assigned=%d Running=%d Succeeded=%d Failed=%d Evicted=%d\n",
		counts[api.Pending],
		counts[api.Assigned],
		counts[api.Running],
		counts[api.Succeeded],
		counts[api.Failed],
		counts[api.Evicted],
	)
}

// printNodeUtilizationLocked prints resource utilization for each node
func (cs *ClusterState) printNodeUtilizationLocked() {
	fmt.Println("\nNode Utilization:")

	for _, node := range cs.Nodes {
		fmt.Printf(
			"  %s | CPU %d/%d | MEM %d/%d\n",
			node.ID,
			node.Allocated.CPU,
			node.TotalCapacity.CPU,
			node.Allocated.Memory,
			node.TotalCapacity.Memory,
		)
	}
}

// printNodeHealthLocked prints the health status of each node
func (cs *ClusterState) printNodeHealthLocked() {
	fmt.Println("\nNode Health:")

	for _, node := range cs.Nodes {
		hb := "never"
		if !node.LastHeartbeat.IsZero() {
			hb = node.LastHeartbeat.Format("15:04:05")
		}

		fmt.Printf(
			"  %s | %s | last heartbeat=%s\n",
			node.ID,
			node.Health,
			hb,
		)
	}
}

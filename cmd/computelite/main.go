package main

import (
	"github.com/adnant1/computelite/pkg/cluster"
	"github.com/adnant1/computelite/pkg/resource"
	"github.com/adnant1/computelite/pkg/scheduler"
)

// The computelite binary simulates creating a cluster, adding nodes,
// submitting jobs, and scheduling them onto the nodes.
func main() {
	clusterInstance := cluster.NewCluster()
	schedulerInstance := scheduler.NewScheduler(clusterInstance)

	// Register nodes with different capacities
	schedulerInstance.Cluster.AddNode(&cluster.Node{
		ID: "node-small",
		TotalCapacity: resource.Resource{
			CPU: 	2000, // 2 cores
			Memory: 4096, // 4 GB
		},
	})

	schedulerInstance.Cluster.AddNode(&cluster.Node{
		ID: "node-medium",
		TotalCapacity: resource.Resource{
			CPU: 	4000, // 4 cores
			Memory: 8192, // 8 GB
		},
	})

	schedulerInstance.Cluster.AddNode(&cluster.Node{
		ID: "node-large",
		TotalCapacity: resource.Resource{
			CPU: 	8000, // 8 cores
			Memory: 16384, // 16 GB
		},
	})

	// Submit jobs with varying resource requirements
	schedulerInstance.SubmitJob(&scheduler.Job{
		ID: 1,
		Requires: resource.Resource{
			CPU:    500,   // 0.5 core
			Memory: 1024,  // 1 GB
		},
	})

	schedulerInstance.SubmitJob(&scheduler.Job{
		ID: 2,
		Requires: resource.Resource{
			CPU:    1000,  // 1 core
			Memory: 2048,  // 2 GB
		},
	})

	schedulerInstance.SubmitJob(&scheduler.Job{
		ID: 3,
		Requires: resource.Resource{
			CPU:    3000,  // 3 cores
			Memory: 4096,  // 4 GB
		},
	})

	schedulerInstance.SubmitJob(&scheduler.Job{
		ID: 4,
		Requires: resource.Resource{
			CPU:    6000,  // 6 cores
			Memory: 8192,  // 8 GB
		},
	})

	schedulerInstance.SubmitJob(&scheduler.Job{
		ID: 5,
		Requires: resource.Resource{
			CPU:    9000,  // 9 cores (will NOT fit anywhere)
			Memory: 4096,
		},
	})

	schedulerInstance.ScheduleAll()

	// Print cluster and scheduling state
	clusterInstance.PrintSummary()
	clusterInstance.PrintNodeUtilization()
	clusterInstance.PrintRunningJobs()
	schedulerInstance.PrintPendingJobs()
}
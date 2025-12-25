package main

import (
	"github.com/adnant1/computelite/pkg/api"
	"github.com/adnant1/computelite/pkg/cluster"
	"github.com/adnant1/computelite/pkg/scheduler"
)

// The computelite binary simulates creating a cluster, adding nodes,
// submitting jobs, and scheduling them onto the nodes.
func main() {
	clusterInstance := cluster.NewCluster()
	schedulerInstance := scheduler.NewScheduler(clusterInstance)

	// Register nodes with different capacities
	schedulerInstance.Cluster.AddNode(&api.Node{
		ID: "node-small",
		TotalCapacity: api.Resource{
			CPU: 	2000, // 2 cores
			Memory: 4096, // 4 GB
		},
	})

	schedulerInstance.Cluster.AddNode(&api.Node{
		ID: "node-medium",
		TotalCapacity: api.Resource{
			CPU: 	4000, // 4 cores
			Memory: 8192, // 8 GB
		},
	})

	schedulerInstance.Cluster.AddNode(&api.Node{
		ID: "node-large",
		TotalCapacity: api.Resource{
			CPU: 	8000, // 8 cores
			Memory: 16384, // 16 GB
		},
	})

	// Submit jobs with varying resource requirements
	schedulerInstance.SubmitJob(&api.Job{
		ID: 1,
		Requires: api.Resource{
			CPU:    500,   // 0.5 core
			Memory: 1024,  // 1 GB
		},
	})

	schedulerInstance.SubmitJob(&api.Job{
		ID: 2,
		Requires: api.Resource{
			CPU:    1000,  // 1 core
			Memory: 2048,  // 2 GB
		},
	})

	schedulerInstance.SubmitJob(&api.Job{
		ID: 3,
		Requires: api.Resource{
			CPU:    3000,  // 3 cores
			Memory: 4096,  // 4 GB
		},
	})

	schedulerInstance.SubmitJob(&api.Job{
		ID: 4,
		Requires: api.Resource{
			CPU:    6000,  // 6 cores
			Memory: 8192,  // 8 GB
		},
	})

	schedulerInstance.SubmitJob(&api.Job{
		ID: 5,
		Requires: api.Resource{
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
package app

import (
	"log"
	"time"

	"github.com/adnant1/computelite/pkg/api"
	"github.com/adnant1/computelite/pkg/cluster"
)

// runScenario bootstraps a simple demo workload
func runScenario(cs *cluster.ClusterState , stop <-chan struct{}) {
	log.Println("[scenario] starting basic scenario")

	addNodes(cs)

	// Start node agents for each node
	go startNodeAgent(cs, "node-small", 1*time.Second, stop)
	go startNodeAgent(cs, "node-medium", 1*time.Second, stop)
	go startNodeAgent(cs, "node-large", 1*time.Second, stop)

	submitJobsOverTime(cs)

	// Optional: simulate a node failure after some time
	go simulateNodeFailure(cs)

	log.Println("[scenario] basic scenario initialized")
}

// addNodes adds a set of nodes with varying capacities to the cluster
func addNodes(cs *cluster.ClusterState) {
	cs.AddNode(&api.Node{
		ID: "node-small",
		TotalCapacity: api.Resource{
			CPU:    2000,
			Memory: 4096,
		},
		Health: api.Healthy,
	})

	cs.AddNode(&api.Node{
		ID: "node-medium",
		TotalCapacity: api.Resource{
			CPU:    4000,
			Memory: 8192,
		},
		Health: api.Healthy,
	})

	cs.AddNode(&api.Node{
		ID: "node-large",
		TotalCapacity: api.Resource{
			CPU:    8000,
			Memory: 16384,
		},
		Health: api.Healthy,
	})
}

// submitJobsOverTime submits jobs to the cluster at regular intervals
func submitJobsOverTime(cs *cluster.ClusterState) {
	go func() {
		for i := int64(1); i <= 6; i++ {
			job := &api.Job{
				ID: i,
				Requires: api.Resource{
					CPU:    1000,
					Memory: 1024,
				},
			}

			if err := cs.SubmitJob(job); err != nil {
				log.Printf("[scenario] failed to submit job=%d: %v", i, err)
			}

			time.Sleep(700 * time.Millisecond)
		}
	}()
}

// simulateNodeFailure marks a node as unhealthy after a delay
func simulateNodeFailure(cs *cluster.ClusterState) {
	time.Sleep(5 * time.Second)

	log.Println("[scenario] simulating node failure: node-medium")
	if err := cs.UpdateNodeHealth("node-medium", api.Unhealthy); err != nil {
		log.Printf("[scenario] failed to mark node unhealthy: %v", err)
	}
}



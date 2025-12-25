package cluster

import (
	"fmt"
	"log"
	"time"

	"github.com/adnant1/computelite/pkg/api"
)

type ClusterState struct {
	Nodes		map[string]*api.Node // Mapping of node IDs to Node structs
	Jobs		map[int64]*api.Job   // Mapping of job IDs to their Job structs
}

// NewCluster initializes and returns a new ClusterState
func NewCluster() *ClusterState {
	return &ClusterState{
		Nodes: make(map[string]*api.Node),
		Jobs:  make(map[int64]*api.Job),
	}
}

// AddNode adds a new node to the cluster state
func (cs *ClusterState) AddNode(node *api.Node) {
	cs.Nodes[node.ID] = node
	log.Printf("[cluster] node=%s added (cpu=%d, memory=%d)\n",
		node.ID,
		node.TotalCapacity.CPU,
		node.TotalCapacity.Memory,
	)
}

// RecordHeartbeat updates the last heartbeat timestamp for a given node
func (cs *ClusterState) RecordHeartbeat(nodeID string) error {
	node, exists := cs.Nodes[nodeID]
	if !exists {
		return fmt.Errorf("node with ID %s not found", nodeID)
	}

	node.LastHeartbeat = time.Now()
	return nil
}

// UpdateNodeHealth updates the health status of a given node
func (cs *ClusterState) UpdateNodeHealth(nodeID string, newHealth api.NodeHealth) error {
	node, exists := cs.Nodes[nodeID]
	if !exists {
		return fmt.Errorf("node with ID %s not found", nodeID)
	}

	if node.Health != newHealth {
		node.Health = newHealth
	}

	return nil
}

// UpdateJobState updates the state of a given job enforcing valid state transitions
func (cs *ClusterState) UpdateJobState(jobID int64, newState api.JobState) error {
	job, exists := cs.Jobs[jobID]
	if !exists {
		return fmt.Errorf("job with ID %d not found", jobID)
	}

	if job.State == newState {
		return nil
	}

	allowed := map[api.JobState]map[api.JobState]bool{
		api.Pending:   {api.Assigned: true},
		api.Assigned:  {api.Running: true, api.Evicted: true},
		api.Running:   {api.Succeeded: true, api.Failed: true, api.Evicted: true},
		api.Evicted:   {api.Pending: true},
		api.Succeeded: {},
		api.Failed:    {},
	}

	if !allowed[job.State][newState] {
		return fmt.Errorf("job %d: invalid state transition from %v to %v", job.ID, job.State, newState)
	}

	job.State = newState
	return nil
}


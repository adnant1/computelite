package cluster

import (
	"fmt"
	"log"
	"time"

	"github.com/adnant1/computelite/pkg/api"
)

type ClusterState struct {
	Nodes		map[string]*api.Node // Mapping of node IDs to Node structs
	RunningJobs	map[int64]*api.Node  // Mapping of running job IDs to their assigned nodes
}

// NewCluster initializes and returns a new ClusterState
func NewCluster() *ClusterState {
	return &ClusterState{
		Nodes:		 make(map[string]*api.Node),
		RunningJobs: make(map[int64]*api.Node),
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
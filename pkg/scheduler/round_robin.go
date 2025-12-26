package scheduler

import (
	"sort"

	"github.com/adnant1/computelite/pkg/api"
)

// RoundRobinPolicy implements the round-robin scheduling policy
type RoundRobinPolicy struct {
	lastIdx int
}

// SelectNode selects the next node in a round-robin fashion
func (p *RoundRobinPolicy) SelectNode(job *api.Job, nodes map[string]*api.Node) (string, bool) {
	nodeIDs := make([]string, 0, len(nodes))
	for nodeID := range nodes {
		nodeIDs = append(nodeIDs, nodeID)
	}
	sort.Strings(nodeIDs)

	n := len(nodeIDs)
	if n == 0 {
		return "", false
	}

	// If the index is invalid, reset it
	if p.lastIdx >= n {
		p.lastIdx = -1
	}

	startIdx := (p.lastIdx + 1) % n
	for i := 0; i < n; i++ {
		idx := (startIdx + i) % n
		nodeID := nodeIDs[idx]
		node := nodes[nodeID]

		availableCPU := node.TotalCapacity.CPU - node.Allocated.CPU
		availableMemory := node.TotalCapacity.Memory - node.Allocated.Memory

		if availableCPU >= job.Requires.CPU && availableMemory >= job.Requires.Memory {
			p.lastIdx = idx
			return nodeID, true
		}
	}

	return "", false
}
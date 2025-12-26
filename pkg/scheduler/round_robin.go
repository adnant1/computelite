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
func (p *RoundRobinPolicy) SelectNode(job *api.Job, nodes []*api.Node) (string, bool) {
	n := len(nodes)
	if n == 0 {
		return "", false
	}

	// Sort nodes by ID for consistent ordering
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].ID < nodes[j].ID
	})

	// If the index is invalid, reset it
	if p.lastIdx >= n {
		p.lastIdx = -1
	}

	startIdx := (p.lastIdx + 1) % n
	for i := 0; i < n; i++ {
		idx := (startIdx + i) % n
		node := nodes[idx]

		availableCPU := node.TotalCapacity.CPU - node.Allocated.CPU
		availableMemory := node.TotalCapacity.Memory - node.Allocated.Memory

		if availableCPU >= job.Requires.CPU && availableMemory >= job.Requires.Memory {
			p.lastIdx = idx
			return node.ID, true
		}
	}

	return "", false
}
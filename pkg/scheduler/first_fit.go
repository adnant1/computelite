package scheduler

import "github.com/adnant1/computelite/pkg/api"

// FirstFitPolicy implements the first-fit scheduling policy
type FirstFitPolicy struct{}

// SelectNode selects the first node that has enough resources for the job
func (p *FirstFitPolicy) SelectNode(job *api.Job, nodes map[string]*api.Node) (string, bool) {
	for nodeID, node := range nodes {
		availableCPU := node.TotalCapacity.CPU - node.Allocated.CPU
		availableMemory := node.TotalCapacity.Memory - node.Allocated.Memory

		if availableCPU >= job.Requires.CPU && availableMemory >= job.Requires.Memory {
			return nodeID, true
		}
	}

	return "", false
}
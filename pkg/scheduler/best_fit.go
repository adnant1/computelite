package scheduler

import "github.com/adnant1/computelite/pkg/api"

// BestFitPolicy implements the best-fit scheduling policy
type BestFitPolicy struct{}

// SelectNode selects the node that best fits the job's resource requirements
func (p *BestFitPolicy) SelectNode(job *api.Job, nodes map[string]*api.Node) (string, bool) {
	var bestNodeID string
	var bestScore int64 = -1

	for nodeID, node := range nodes {
		availableCPU := node.TotalCapacity.CPU - node.Allocated.CPU
		availableMemory := node.TotalCapacity.Memory - node.Allocated.Memory

		if availableCPU >= job.Requires.CPU && availableMemory >= job.Requires.Memory {
			cpuDiff := availableCPU - job.Requires.CPU
			memDiff := availableMemory - job.Requires.Memory
			score := cpuDiff + memDiff

			if bestScore == -1 || score < bestScore {
				bestScore = score
				bestNodeID = nodeID
			}
		}
	}

	return bestNodeID, bestScore != -1
}

package cluster

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/adnant1/computelite/pkg/api"
)

type ClusterState struct {
	Nodes		map[string]*api.Node // Mapping of node IDs to Node structs
	Jobs		map[int64]*api.Job   // Mapping of job IDs to their Job structs

	mu 	  		sync.RWMutex
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
	cs.mu.Lock()
	defer cs.mu.Unlock()

	cs.Nodes[node.ID] = node
	log.Printf("[cluster] node=%s added (cpu=%d, memory=%d)",
		node.ID,
		node.TotalCapacity.CPU,
		node.TotalCapacity.Memory,
	)
}

// SubmitJob adds a new job to the cluster state
func (cs *ClusterState) SubmitJob(job *api.Job) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	_, exists := cs.Jobs[job.ID]
	if exists {
		return fmt.Errorf("job with ID %d already exists", job.ID)
	}

	cs.Jobs[job.ID] = job
	err := cs.updateJobStateLocked(job.ID, api.Pending)
	if err != nil {
		return err
	}

	log.Printf("[cluster] job=%d submitted (cpu=%d, mem=%d)",
		job.ID,
		job.Requires.CPU,
		job.Requires.Memory,
	)
	return nil
}

// RecordHeartbeat updates the last heartbeat timestamp for a given node
func (cs *ClusterState) RecordHeartbeat(nodeID string) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	node, exists := cs.Nodes[nodeID]
	if !exists {
		return fmt.Errorf("node with ID %s not found", nodeID)
	}

	node.LastHeartbeat = time.Now()
	return nil
}

// UpdateNodeHealth updates the health status of a given node
func (cs *ClusterState) UpdateNodeHealth(nodeID string, newHealth api.NodeHealth) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

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
	cs.mu.Lock()
	defer cs.mu.Unlock()

	return cs.updateJobStateLocked(jobID, newState)
}

// updateJobStateLocked is the internal version that assumes the lock is already held
func (cs *ClusterState) updateJobStateLocked(jobID int64, newState api.JobState) error {
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

// EvictAndRequeueJob evicts a job from its assigned node and re-queues it
func (cs *ClusterState) EvictAndRequeueJob(jobID int64) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	job, exists := cs.Jobs[jobID]
	if !exists {
		return fmt.Errorf("job with ID %d not found", jobID)
	}

	if job.State == api.Succeeded || job.State == api.Failed {
		return fmt.Errorf("job %d is already completed with state %v", jobID, job.State)
	}

	if job.AssignedNodeID == "" && job.State == api.Pending {
		return nil
	}

	node, exists := cs.Nodes[job.AssignedNodeID]
	if !exists {
		return fmt.Errorf("node with ID %s not found", job.AssignedNodeID)
	}

	// Free up resources on the node
	node.Allocated.CPU -= job.Requires.CPU
	node.Allocated.Memory -= job.Requires.Memory

	// Safe guards against negative allocation
	if node.Allocated.CPU < 0 {
		node.Allocated.CPU = 0
	}

	if node.Allocated.Memory < 0 {
		node.Allocated.Memory = 0
	}

	job.AssignedNodeID = ""

	err := cs.updateJobStateLocked(jobID, api.Evicted)
	if err != nil {
		return err
	}

	err = cs.updateJobStateLocked(jobID, api.Pending)
	if err != nil {
		return err
	}

	return nil
}

// AssignJob assigns a pending job to a healthy node if resources permit
func (cs *ClusterState) AssignJob(jobID int64, nodeID string) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	job, jobExists := cs.Jobs[jobID]
	if !jobExists {
		return fmt.Errorf("job with ID %d not found", jobID)
	}

	node, nodeExists := cs.Nodes[nodeID]
	if !nodeExists {
		return fmt.Errorf("node with ID %s not found", nodeID)
	}

	if job.State != api.Pending {
		return fmt.Errorf("job %d is not in pending state", jobID)
	}

	if node.Health != api.Healthy {
		return fmt.Errorf("node %s is not healthy", nodeID)
	}

	if node.Allocated.CPU + job.Requires.CPU > node.TotalCapacity.CPU ||
		node.Allocated.Memory + job.Requires.Memory > node.TotalCapacity.Memory {
		return fmt.Errorf("node %s does not have enough resources for job %d", nodeID, jobID)
	}

	err := cs.updateJobStateLocked(jobID, api.Assigned)
	if err != nil {
		return err
	}

	node.Allocated.CPU += job.Requires.CPU
	node.Allocated.Memory += job.Requires.Memory

	job.AssignedNodeID = nodeID

	return nil
}

// ListJobs returns a slice of all jobs in the cluster
func (cs *ClusterState) ListJobs() []*api.Job {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	jobs := make([]*api.Job, 0, len(cs.Jobs))
	for _, job := range cs.Jobs {
		jobs = append(jobs, job)
	}

	return jobs
}

// ListJobsByState returns a slice of jobs filtered by the specified state
func (cs *ClusterState) ListJobsByState(state api.JobState) []*api.Job {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	jobs := make([]*api.Job, 0)
	for _, job := range cs.Jobs {
		if job.State == state {
			jobs = append(jobs, job)
		}
	}

	return jobs
}

// ListNodes returns a slice of all nodes in the cluster
func (cs *ClusterState) ListNodes() []*api.Node {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	nodes := make([]*api.Node, 0, len(cs.Nodes))
	for _, node := range cs.Nodes {
		nodes = append(nodes, node)
	}

	return nodes
}

// ListHealthyNodes returns a slice of all healthy nodes in the cluster
func (cs *ClusterState) ListHealthyNodes() []*api.Node {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	nodes := make([]*api.Node, 0)
	for _, node := range cs.Nodes {
		if node.Health == api.Healthy {
			nodes = append(nodes, node)
		}
	}

	return nodes
}
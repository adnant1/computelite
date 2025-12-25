package api

import "time"

type NodeHealth int

// Define node health states
const (
	Unknown NodeHealth = iota
	Healthy
	Unhealthy
)

type Node struct {
	ID            string

	TotalCapacity Resource // Total resource capacity of the node
	Allocated     Resource // Currently allocated resources on the node

	Health        NodeHealth	    // Current health status of the node
	LastHeartbeat time.Time         // Last heartbeat timestamp from the node
}

// String returns a string representation of the NodeHealth
func (h NodeHealth) String() string {
	switch h {
	case Healthy:
		return "Healthy"
	case Unhealthy:
		return "Unhealthy"
	default:
		return "Unknown"
	}
}
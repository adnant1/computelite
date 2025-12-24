package cluster

import (
	"time"

	"github.com/adnant1/computelite/pkg/resource"
)

type NodeHealth int

// Define node health states
const (
	Unknown NodeHealth = iota
	Healthy
	Unhealthy
)

type Node struct {
	ID            string

	TotalCapacity resource.Resource // Total resource capacity of the node
	Allocated     resource.Resource // Currently allocated resources on the node

	Health        NodeHealth	    // Current health status of the node
	LastHeartbeat time.Time         // Last heartbeat timestamp from the node
}
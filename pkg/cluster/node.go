package cluster

import "github.com/adnant1/computelite/pkg/resource"

type Node struct {
	ID            string            
	TotalCapacity resource.Resource // Total resource capacity of the node
	Allocated     resource.Resource // Currently allocated resources on the node
}
package model

type Node struct {
	ID            string   // Unique identifier for the node
	TotalCapacity Resource // Total resource capacity of the node
	Allocated     Resource // Currently allocated resources on the node
}
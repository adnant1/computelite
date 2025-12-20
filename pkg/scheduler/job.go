package scheduler

import "github.com/adnant1/computelite/pkg/resource"

type JobState int

// Define job states
const (
	Pending JobState = iota
	Running
)

type Job struct {
	ID       int64
	Requires resource.Resource // Resources required by the job
	State    JobState          // Current state of the job
}
package scheduler

import "github.com/adnant1/computelite/pkg/api"

// SchedulerPolicy defines the interface for scheduling policies
type SchedulerPolicy interface {
	SelectNode(job *api.Job, nodes []*api.Node) (string, bool)
}
package app

import (
	"time"

	"github.com/adnant1/computelite/pkg/cluster"
	"github.com/adnant1/computelite/pkg/controller"
	"github.com/adnant1/computelite/pkg/scheduler"
)

// startControllers initializes and starts all controllers for the cluster state
func startControllers(cs *cluster.ClusterState, stop <-chan struct{}) {
	policy := &scheduler.BestFitPolicy{}

	schedulerController := controller.NewSchedulerController(
		cs,
		policy,
		300*time.Millisecond,
	)

	jobController := controller.NewJobController(
		cs,
		500*time.Millisecond,
	)

	healthController := controller.NewHealthController(
		cs,
		1*time.Second,
	)

	reschedulerController := controller.NewReschedulerController(
		cs,
		1*time.Second,
	)

	go schedulerController.Run(stop)
	go jobController.Run(stop)
	go healthController.Run(stop)
	go reschedulerController.Run(stop)
}

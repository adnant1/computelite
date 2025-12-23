package scheduler

import "fmt"

// PrintPendingJobs prints all pending jobs in the scheduler queue.
func (s *Scheduler) PrintPendingJobs() {
	fmt.Println("\n=== Pending Jobs ===")

	jobs := s.PendingJobs.Jobs()
	if len(jobs) == 0 {
		fmt.Println("No pending jobs")
		return
	}

	for _, job := range jobs {
		fmt.Printf("Job %d (CPU=%d, MEM=%d)\n",
			job.ID,
			job.Requires.CPU,
			job.Requires.Memory,
		)
	}
}
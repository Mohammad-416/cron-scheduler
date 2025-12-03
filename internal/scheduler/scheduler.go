package scheduler

import (
	"context"
	"fmt"
	"os/exec"
	"sync"
	"time"

	"cron-scheduler/internal/config"
)

type Scheduler struct {
	jobs  []config.JobConfig
	state *StateManager
	queue chan config.JobConfig
	wg    sync.WaitGroup
}

func NewScheduler(jobs []config.JobConfig, state *StateManager) *Scheduler {
	return &Scheduler{
		jobs:  jobs,
		state: state,
		queue: make(chan config.JobConfig, 100),
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	for _, job := range s.jobs {
		j := job // copy
		s.wg.Add(1)
		go s.jobLoop(ctx, j)
	}

	// single worker for simplicity — you can extend later
	s.wg.Add(1)
	go s.worker(ctx)
}

func (s *Scheduler) jobLoop(ctx context.Context, job config.JobConfig) {
	defer s.wg.Done()

	interval := job.EveryDuration
	now := time.Now()

	// Load last run time
	s.state.mu.RLock()
	lastStatus := s.state.Data[job.Name]
	s.state.mu.RUnlock()

	runNow, wait := computeNextRun(lastStatus.LastRun, interval, now)

	if runNow {
		fmt.Println("[resume/run] Running immediately:", job.Name)
		s.enqueue(ctx, job)
	}

	timer := time.NewTimer(wait)
	defer timer.Stop()

	fmt.Printf("[scheduler] %s next run in %v\n", job.Name, wait)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("[scheduler] exiting job loop:", job.Name)
			return

		case <-timer.C:
			// Run job
			s.enqueue(ctx, job)

			// Reset for next interval
			timer.Reset(interval)
		}
	}
}

func (s *Scheduler) enqueue(ctx context.Context, job config.JobConfig) {
	select {
	case s.queue <- job:
	default:
		fmt.Println("[error] queue is full, dropping job:", job.Name)
	}
}

func (s *Scheduler) worker(ctx context.Context) {
	defer s.wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("[worker] shutting down")
			return

		case job := <-s.queue:
			s.runJob(ctx, job)
		}
	}
}

func (s *Scheduler) runJob(ctx context.Context, job config.JobConfig) {
	fmt.Println("[run]", job.Name, "→", job.Command)

	cmd := exec.CommandContext(ctx, "sh", "-c", job.Command)
	out, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("[error]", job.Name, ":", err)
	} else {
		fmt.Println("[output]", job.Name, ":", string(out))
	}

	// Update last run
	s.state.mu.Lock()
	st := s.state.Data[job.Name]
	st.LastRun = time.Now()
	st.RunCount++
	s.state.Data[job.Name] = st
	s.state.mu.Unlock()

	// Save state after each run
	s.state.Save()
}

func (s *Scheduler) Wait() {
	s.wg.Wait()
}

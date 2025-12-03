package types

import "time"

type JobStatus struct {
	LastRun  time.Time `json:"last_run"`
	RunCount int       `json:"run_count"`
}

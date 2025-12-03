package scheduler

import (
	"time"
)

func computeNextRun(lastRun time.Time, interval time.Duration, now time.Time) (runNow bool, wait time.Duration) {
	if lastRun.IsZero() {
		// Never run before
		return true, interval
	}

	elapsed := now.Sub(lastRun)

	if elapsed >= interval {
		// Missed at least one full interval → run immediately
		remaining := interval - (elapsed % interval)
		return true, remaining
	}

	// Paused before interval completed → resume leftover time
	remaining := interval - elapsed
	return false, remaining
}

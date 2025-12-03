package scheduler

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"cron-scheduler/internal/types"
)

type StateManager struct {
	path string
	mu   sync.RWMutex
	Data map[string]types.JobStatus
}

func NewStateManager(path string) *StateManager {
	return &StateManager{
		path: path,
		Data: make(map[string]types.JobStatus),
	}
}

func (sm *StateManager) Load() {
	data, err := os.ReadFile(sm.path)
	if err != nil {
		fmt.Println("[state] no state file found:", err)
		return
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()

	if err := json.Unmarshal(data, &sm.Data); err != nil {
		fmt.Println("[state] error unmarshaling:", err)
	}
	fmt.Println("[state] loaded job run history")
}

func (sm *StateManager) Save() {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	data, err := json.MarshalIndent(sm.Data, "", "  ")
	if err != nil {
		fmt.Println("[state] marshal error:", err)
		return
	}

	if err := os.WriteFile(sm.path, data, 0644); err != nil {
		fmt.Println("[state] write error:", err)
	}
}

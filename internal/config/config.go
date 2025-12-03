package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type JobConfig struct {
	Name    string `yaml:"name"`
	Every   string `yaml:"every"`
	Command string `yaml:"command"`

	// parsed duration (not from YAML)
	EveryDuration time.Duration `yaml:"-"`
}

type Config struct {
	Jobs []JobConfig `yaml:"jobs"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal yaml: %w", err)
	}

	for i := range cfg.Jobs {
		j := &cfg.Jobs[i]

		if j.Name == "" {
			j.Name = fmt.Sprintf("job_%d", i)
		}
		if j.Every == "" {
			return nil, fmt.Errorf("job %q missing 'every' field", j.Name)
		}
		if j.Command == "" {
			return nil, fmt.Errorf("job %q missing 'command' field", j.Name)
		}

		d, err := time.ParseDuration(j.Every)
		if err != nil {
			return nil, fmt.Errorf("job %q invalid duration %q: %w", j.Name, j.Every, err)
		}
		j.EveryDuration = d
	}

	return &cfg, nil
}

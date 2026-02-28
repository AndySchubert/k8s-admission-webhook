package admission

import (
	"os"

	"gopkg.in/yaml.v3"
)

// PolicyConfig controls which admission policies are enforced.
type PolicyConfig struct {
	Policies struct {
		DenyLatestTag    bool `yaml:"denyLatestTag"`
		RequireResources bool `yaml:"requireResources"`
	} `yaml:"policies"`
}

// DefaultConfig returns a PolicyConfig with all policies enabled.
func DefaultConfig() *PolicyConfig {
	cfg := &PolicyConfig{}
	cfg.Policies.DenyLatestTag = true
	cfg.Policies.RequireResources = true
	return cfg
}

// LoadConfig reads a YAML policy file from disk.
// Returns DefaultConfig if the file does not exist.
func LoadConfig(path string) (*PolicyConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, err
	}

	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

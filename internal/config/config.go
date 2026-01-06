package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration
type Config struct {
	Model          string       `yaml:"model"`
	Language       string       `yaml:"language"`
	Types          []CommitType `yaml:"types"`
	Template       string       `yaml:"template"`
	Scopes         []string     `yaml:"scopes"`
	CustomPrompt   string       `yaml:"custom_prompt,omitempty"`
	MaxDiffLength  int          `yaml:"max_diff_length,omitempty"`
	DetailedCommit bool         `yaml:"detailed_commit,omitempty"` // Generate detailed commit messages with body
}

// CommitType defines a type of commit with description and emoji
type CommitType struct {
	Name  string `yaml:"name"`
	Desc  string `yaml:"desc"`
	Emoji string `yaml:"emoji"`
}

// LoadConfig loads configuration from file with fallback to defaults
func LoadConfig() (*Config, error) {
	// Try to find config file in order of priority
	configPaths := []string{
		".gitcommit.yaml",
		".gitcommit.yml",
		filepath.Join(os.Getenv("HOME"), ".gitcommit.yaml"),
		filepath.Join(os.Getenv("HOME"), ".gitcommit.yml"),
	}

	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			// Found config file
			config, err := loadFromFile(path)
			if err != nil {
				return nil, fmt.Errorf("invalid config file format at %s: %w\nCheck .gitcommit.yaml syntax", path, err)
			}
			return config, nil
		}
	}

	// No config file found, use defaults
	return DefaultConfig(), nil
}

// loadFromFile loads configuration from a YAML file
func loadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	// Apply defaults for missing fields
	if config.Model == "" {
		config.Model = "qwen2.5-coder:7b"
	}
	if config.Language == "" {
		config.Language = "en"
	}
	if len(config.Types) == 0 {
		config.Types = DefaultConfig().Types
	}
	if config.Template == "" {
		config.Template = "{type}{scope}: {emoji} {message}"
	}
	if config.MaxDiffLength == 0 {
		config.MaxDiffLength = 2000
	}

	return config, nil
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Model:    "qwen2.5-coder:7b",
		Language: "en",
		Types: []CommitType{
			{Name: "feat", Desc: "A new feature", Emoji: "✨"},
			{Name: "fix", Desc: "A bug fix", Emoji: "🐛"},
			{Name: "docs", Desc: "Documentation only changes", Emoji: "📝"},
			{Name: "style", Desc: "Code style changes (formatting, etc)", Emoji: "💄"},
			{Name: "refactor", Desc: "Code refactoring", Emoji: "♻️"},
			{Name: "perf", Desc: "Performance improvements", Emoji: "⚡"},
			{Name: "test", Desc: "Adding or updating tests", Emoji: "✅"},
			{Name: "chore", Desc: "Build process or auxiliary tool changes", Emoji: "🔧"},
			{Name: "ci", Desc: "CI configuration changes", Emoji: "👷"},
			{Name: "build", Desc: "Build system changes", Emoji: "📦"},
		},
		Template:       "{type}{scope}: {emoji} {message}",
		Scopes:         []string{},
		MaxDiffLength:  2000,
		DetailedCommit: true, // Default to detailed commits
	}
}

// GetTypeByName finds a commit type by its name
func (c *Config) GetTypeByName(name string) *CommitType {
	for _, t := range c.Types {
		if t.Name == name {
			return &t
		}
	}
	return nil
}

// Save saves the configuration to a file
func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

package src

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the complete configuration structure
type Config struct {
	API       APIConfig      `yaml:"api"`
	Languages LanguageConfig `yaml:"languages"`
	Cache     CacheConfig    `yaml:"cache"`
	Advanced  AdvancedConfig `yaml:"advanced"`
}

// APIConfig represents API-related configuration
type APIConfig struct {
	Endpoint    string  `yaml:"endpoint"`
	Key         string  `yaml:"key"`
	Model       string  `yaml:"model"`
	Timeout     int     `yaml:"timeout"`
	MaxTokens   int     `yaml:"max_tokens"`
	Temperature float64 `yaml:"temperature"`
}

// LanguageConfig represents language-related configuration
type LanguageConfig struct {
	Common   []string `yaml:"common"`
	Priority []string `yaml:"priority"`
}

// CacheConfig represents cache-related configuration
type CacheConfig struct {
	Enabled   bool   `yaml:"enabled"`
	Directory string `yaml:"directory"`
	TTL       int    `yaml:"ttl"`
}

// AdvancedConfig represents advanced configuration options
type AdvancedConfig struct {
	Debug          bool   `yaml:"debug"`
	LogDir         string `yaml:"log_dir"`
	PromptTemplate string `yaml:"prompt_template"`
}

// DefaultConfig returns a Config with default values
func DefaultConfig() *Config {
	return &Config{
		API: APIConfig{
			Endpoint:    "https://api.openai.com/v1/chat/completions",
			Key:         "",
			Model:       "gpt-4",
			Timeout:     30,
			MaxTokens:   1000,
			Temperature: 0.7,
		},
		Languages: LanguageConfig{
			Common:   []string{"zh", "en", "ja", "ko", "es", "fr", "de", "ru", "pt", "it"},
			Priority: []string{"zh", "en"},
		},
		Cache: CacheConfig{
			Enabled:   true,
			Directory: ".cache/fanyi",
			TTL:       720, // 30 days in hours
		},
		Advanced: AdvancedConfig{
			Debug:  false,
			LogDir: ".log/fanyi",
			PromptTemplate: `You are a professional translator. Translate the following text to {language}.
Only return the translated text without any explanation or additional content.

Text: {input_text}`,
		},
	}
}

// Load loads configuration from file, environment variables, and defaults
func Load() (*Config, error) {
	cfg := DefaultConfig()

	// Try to load from config file
	homeDir, err := os.UserHomeDir()
	if err == nil {
		configPath := filepath.Join(homeDir, ".config/fanyi/config.yaml")
		if data, err := os.ReadFile(configPath); err == nil {
			if err := yaml.Unmarshal(data, cfg); err != nil {
				return nil, fmt.Errorf("failed to parse config file: %w", err)
			}
		}
	}

	// Override with environment variables
	cfg.applyEnvVars()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// applyEnvVars applies environment variables to the config
func (c *Config) applyEnvVars() {
	// API configuration
	if endpoint := os.Getenv("FANYI_API_ENDPOINT"); endpoint != "" {
		c.API.Endpoint = endpoint
	}
	if key := os.Getenv("FANYI_API_KEY"); key != "" {
		c.API.Key = key
	}
	if model := os.Getenv("FANYI_API_MODEL"); model != "" {
		c.API.Model = model
	}
	if timeout := os.Getenv("FANYI_API_TIMEOUT"); timeout != "" {
		if val, err := strconv.Atoi(timeout); err == nil {
			c.API.Timeout = val
		}
	}
	if maxTokens := os.Getenv("FANYI_API_MAX_TOKENS"); maxTokens != "" {
		if val, err := strconv.Atoi(maxTokens); err == nil {
			c.API.MaxTokens = val
		}
	}

	// Language configuration
	if langs := os.Getenv("FANYI_LANGUAGES"); langs != "" {
		c.Languages.Common = strings.Split(langs, ",")
	}
	if priority := os.Getenv("FANYI_LANGUAGE_PRIORITY"); priority != "" {
		c.Languages.Priority = strings.Split(priority, ",")
	}

	// Cache configuration
	if enabled := os.Getenv("FANYI_CACHE_ENABLED"); enabled != "" {
		c.Cache.Enabled = enabled == "true"
	}
	if dir := os.Getenv("FANYI_CACHE_DIR"); dir != "" {
		c.Cache.Directory = dir
	}

	// Advanced configuration
	if debug := os.Getenv("FANYI_DEBUG"); debug != "" {
		c.Advanced.Debug = debug == "true"
	}
	if logDir := os.Getenv("FANYI_LOG_DIR"); logDir != "" {
		c.Advanced.LogDir = logDir
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.API.Endpoint == "" {
		return fmt.Errorf("API endpoint is required")
	}
	if c.API.Key == "" {
		return fmt.Errorf("API key is required (set FANYI_API_KEY or add to config file)")
	}
	if c.API.Model == "" {
		return fmt.Errorf("API model is required")
	}
	if len(c.Languages.Priority) == 0 {
		return fmt.Errorf("at least one priority language is required")
	}
	return nil
}

// GetCacheDir returns the absolute path to the cache directory
func (c *Config) GetCacheDir() string {
	if filepath.IsAbs(c.Cache.Directory) {
		return c.Cache.Directory
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return c.Cache.Directory
	}
	return filepath.Join(homeDir, c.Cache.Directory)
}

// GetLogDir returns the absolute path to the log directory
func (c *Config) GetLogDir() string {
	if filepath.IsAbs(c.Advanced.LogDir) {
		return c.Advanced.LogDir
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return c.Advanced.LogDir
	}
	return filepath.Join(homeDir, c.Advanced.LogDir)
}

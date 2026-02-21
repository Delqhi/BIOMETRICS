package config

import (
	"biometrics-cli/internal/state"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Config struct {
	mu        sync.RWMutex
	values    map[string]interface{}
	filePath  string
	watchers  []Watcher
	hotReload bool
}

type Watcher func(key string, value interface{})

type ProviderConfig struct {
	Provider string            `json:"provider"`
	Models   map[string]Model  `json:"models"`
	Options  map[string]string `json:"options,omitempty"`
	Enabled  bool              `json:"enabled"`
	APIKey   string            `json:"api_key,omitempty"`
	BaseURL  string            `json:"base_url,omitempty"`
}

type Model struct {
	ID         string                            `json:"id"`
	Name       string                            `json:"name"`
	Limit      ModelLimit                        `json:"limit"`
	Modalities map[string][]string               `json:"modalities,omitempty"`
	Variants   map[string]map[string]interface{} `json:"variants,omitempty"`
}

type ModelLimit struct {
	Context int `json:"context"`
	Output  int `json:"output"`
}

type OpenCodeConfig struct {
	Provider map[string]ProviderConfig `json:"provider"`
	MCP      map[string]MCPConfig      `json:"mcp,omitempty"`
	Plugin   []string                  `json:"plugin,omitempty"`
}

type MCPConfig struct {
	Type        string            `json:"type"`
	Command     []string          `json:"command,omitempty"`
	URL         string            `json:"url,omitempty"`
	Enabled     bool              `json:"enabled"`
	Environment map[string]string `json:"environment,omitempty"`
}

var GlobalConfig = &Config{
	values:    make(map[string]interface{}),
	hotReload: true,
	watchers:  make([]Watcher, 0),
}

func (c *Config) Load(path string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.filePath = path

	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &c.values); err != nil {
		return err
	}

	state.GlobalState.Log("INFO", fmt.Sprintf("Loaded config from: %s", path))

	if c.hotReload {
		go c.watchFile(path)
	}

	return nil
}

func (c *Config) watchFile(path string) {
	watcher, err := NewWatcher(path)
	if err != nil {
		state.GlobalState.Log("ERROR", fmt.Sprintf("Failed to watch config: %v", err))
		return
	}
	defer watcher.Close()

	for {
		select {
		case <-watcher.Events:
			if err := c.Reload(); err != nil {
				state.GlobalState.Log("ERROR", fmt.Sprintf("Config reload failed: %v", err))
			}
		case <-watcher.Done:
			return
		}
	}
}

func (c *Config) Reload() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.filePath == "" {
		return fmt.Errorf("no config file path set")
	}

	data, err := os.ReadFile(c.filePath)
	if err != nil {
		return err
	}

	newValues := make(map[string]interface{})
	if err := json.Unmarshal(data, &newValues); err != nil {
		return err
	}

	oldValues := c.values
	c.values = newValues

	for key, value := range newValues {
		if oldValues[key] != value {
			c.notifyWatchers(key, value)
		}
	}

	state.GlobalState.Log("INFO", "Config reloaded")
	return nil
}

func (c *Config) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, exists := c.values[key]
	return value, exists
}

func (c *Config) GetString(key string, defaultValue string) string {
	if value, ok := c.Get(key); ok {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return defaultValue
}

func (c *Config) GetInt(key string, defaultValue int) int {
	if value, ok := c.Get(key); ok {
		if num, ok := value.(float64); ok {
			return int(num)
		}
	}
	return defaultValue
}

func (c *Config) GetBool(key string, defaultValue bool) bool {
	if value, ok := c.Get(key); ok {
		if b, ok := value.(bool); ok {
			return b
		}
	}
	return defaultValue
}

func (c *Config) GetMap(key string) map[string]interface{} {
	if value, ok := c.Get(key); ok {
		if m, ok := value.(map[string]interface{}); ok {
			return m
		}
	}
	return nil
}

func (c *Config) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.values[key] = value
	c.notifyWatchers(key, value)
}

func (c *Config) SetAll(values map[string]interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, value := range values {
		c.values[key] = value
		c.notifyWatchers(key, value)
	}
}

func (c *Config) Save() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.filePath == "" {
		return fmt.Errorf("no config file path set")
	}

	data, err := json.MarshalIndent(c.values, "", "  ")
	if err != nil {
		return err
	}

	absPath, _ := filepath.Abs(c.filePath)
	if err := os.WriteFile(absPath, data, 0644); err != nil {
		return err
	}

	state.GlobalState.Log("INFO", fmt.Sprintf("Saved config to: %s", c.filePath))
	return nil
}

func (c *Config) Watch(watcher Watcher) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.watchers = append(c.watchers, watcher)
}

func (c *Config) notifyWatchers(key string, value interface{}) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, watcher := range c.watchers {
		watcher(key, value)
	}
}

type FileWatcher struct {
	Events  chan bool
	Done    chan bool
	watcher *os.File
}

func NewWatcher(path string) (*FileWatcher, error) {
	watcher, err := os.Open(filepath.Dir(path))
	if err != nil {
		return nil, err
	}

	fw := &FileWatcher{
		Events:  make(chan bool, 1),
		Done:    make(chan bool),
		watcher: watcher,
	}

	go func() {
		buf := make([]byte, 1024)
		for {
			select {
			case <-fw.Done:
				return
			default:
				n, err := watcher.Read(buf)
				if err != nil {
					select {
					case <-fw.Done:
						return
					case <-time.After(100 * time.Millisecond):
					}
					continue
				}
				if n > 0 {
					select {
					case fw.Events <- true:
					default:
					}
				}
			}
		}
	}()

	return fw, nil
}

func (fw *FileWatcher) Close() {
	close(fw.Done)
	fw.watcher.Close()
}

func LoadOpenCodeConfig(path string) (*OpenCodeConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config OpenCodeConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func GetProviderConfig(name string) (*ProviderConfig, error) {
	path := os.ExpandEnv("$HOME/.config/opencode/opencode.json")
	config, err := LoadOpenCodeConfig(path)
	if err != nil {
		return nil, err
	}

	provider, exists := config.Provider[name]
	if !exists {
		return nil, fmt.Errorf("provider %s not found", name)
	}

	return &provider, nil
}

func GetModelConfig(providerName, modelName string) (*Model, error) {
	provider, err := GetProviderConfig(providerName)
	if err != nil {
		return nil, err
	}

	model, exists := provider.Models[modelName]
	if !exists {
		return nil, fmt.Errorf("model %s not found in provider %s", modelName, providerName)
	}

	return &model, nil
}

func GetMCPConfig(name string) (*MCPConfig, error) {
	path := os.ExpandEnv("$HOME/.config/opencode/opencode.json")
	config, err := LoadOpenCodeConfig(path)
	if err != nil {
		return nil, err
	}

	mcp, exists := config.MCP[name]
	if !exists {
		return nil, fmt.Errorf("MCP %s not found", name)
	}

	return &mcp, nil
}

func init() {
	configPath := os.ExpandEnv("$HOME/.config/opencode/opencode.json")
	if err := GlobalConfig.Load(configPath); err != nil {
		state.GlobalState.Log("WARN", fmt.Sprintf("Failed to load config: %v", err))
	}
}

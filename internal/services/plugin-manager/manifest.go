package plugin_manager

import (
	"os"
	"sync"

	"github.com/vayload/vayload/pkg/encoding/json5"
)

const (
	PluginsDir = "./plugins"
	CacheDir   = ".cache"
	ConfigFile = "plugins.json5"
	BaseURL    = "http://localhost:8080/api/v1"
)

type PluginsManifest struct {
	mu sync.RWMutex

	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"dev_dependencies"`
}

func (pc *PluginsManifest) Add(id, version string, isDev bool) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	if isDev {
		if pc.DevDependencies == nil {
			pc.DevDependencies = make(map[string]string)
		}
		pc.DevDependencies[id] = version
	} else {
		if pc.Dependencies == nil {
			pc.Dependencies = make(map[string]string)
		}
		pc.Dependencies[id] = version
	}
}

func (pc *PluginsManifest) Remove(id string) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	delete(pc.Dependencies, id)
	delete(pc.DevDependencies, id)
}

func (pc *PluginsManifest) Get(id string) (string, bool) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	if version, ok := pc.Dependencies[id]; ok {
		return version, true
	}
	version, ok := pc.DevDependencies[id]
	return version, ok
}

func (pc *PluginsManifest) Save() error {
	pc.mu.RLock()
	data, err := json5.MarshalIndent(pc, "", "    ")
	pc.mu.RUnlock()
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigFile, data, 0644)
}

func (pc *PluginsManifest) Load() error {
	data, err := os.ReadFile(ConfigFile)
	if err != nil {
		return err
	}

	return json5.Unmarshal(data, pc)
}

func (pc *PluginsManifest) Exists() bool {
	_, err := os.Stat(ConfigFile)
	return err == nil
}

func (pc *PluginsManifest) Empty() bool {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	return len(pc.Dependencies) == 0 && len(pc.DevDependencies) == 0
}

func (pc *PluginsManifest) GetDependencies() map[string]string {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	return pc.Dependencies
}

func (pc *PluginsManifest) GetDevDependencies() map[string]string {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	return pc.DevDependencies
}

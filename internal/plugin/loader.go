package plugin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/barun-bash/human/internal/codegen"
)

// PluginManifest describes an installed external plugin.
type PluginManifest struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Category    string `json:"category"` // "frontend", "backend", "database", "infra"
	Binary      string `json:"binary"`
	Source      string `json:"source,omitempty"`      // go module path or URL
	InstalledAt string `json:"installed_at,omitempty"` // RFC3339 timestamp
}

// PluginsDir returns the directory where external plugins are stored.
// Defaults to ~/.human/plugins/.
func PluginsDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", ".human", "plugins")
	}
	return filepath.Join(home, ".human", "plugins")
}

// Discover scans the plugins directory and returns manifests for all
// installed plugins. Plugins with invalid manifests are silently skipped.
func Discover() ([]PluginManifest, error) {
	dir := PluginsDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading plugins directory: %w", err)
	}

	var manifests []PluginManifest
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		m, err := LoadManifest(filepath.Join(dir, entry.Name()))
		if err != nil {
			continue // skip broken plugins
		}
		manifests = append(manifests, m)
	}
	return manifests, nil
}

// LoadManifest reads a plugin.json manifest from the given plugin directory.
func LoadManifest(pluginDir string) (PluginManifest, error) {
	path := filepath.Join(pluginDir, "plugin.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return PluginManifest{}, fmt.Errorf("reading manifest: %w", err)
	}

	var m PluginManifest
	if err := json.Unmarshal(data, &m); err != nil {
		return PluginManifest{}, fmt.Errorf("parsing manifest: %w", err)
	}

	if m.Name == "" {
		return PluginManifest{}, fmt.Errorf("manifest missing name field")
	}
	return m, nil
}

// LoadAll discovers installed plugins and returns them as CodeGenerator
// adapters ready for registration. Plugins that fail to load are skipped.
func LoadAll() ([]codegen.CodeGenerator, error) {
	manifests, err := Discover()
	if err != nil {
		return nil, err
	}

	var generators []codegen.CodeGenerator
	for _, m := range manifests {
		bin := resolveBinary(m)
		if bin == "" {
			continue // binary not found, skip
		}
		generators = append(generators, &ExternalGenerator{
			manifest: m,
			binary:   bin,
		})
	}
	return generators, nil
}

// resolveBinary finds the plugin binary. It first checks the plugin directory,
// then falls back to PATH lookup. On Windows, it also checks for .exe suffix.
func resolveBinary(m PluginManifest) string {
	pluginDir := filepath.Join(PluginsDir(), m.Name)

	// Check plugin directory first.
	candidates := []string{m.Binary}
	if runtime.GOOS == "windows" && filepath.Ext(m.Binary) == "" {
		candidates = append(candidates, m.Binary+".exe")
	}

	for _, name := range candidates {
		p := filepath.Join(pluginDir, name)
		if info, err := os.Stat(p); err == nil && !info.IsDir() {
			return p
		}
	}

	// Fall back to PATH lookup.
	for _, name := range candidates {
		if p, err := findInPath(name); err == nil {
			return p
		}
	}

	return ""
}

// findInPath searches PATH for the given binary name.
func findInPath(name string) (string, error) {
	pathEnv := os.Getenv("PATH")
	for _, dir := range filepath.SplitList(pathEnv) {
		p := filepath.Join(dir, name)
		if info, err := os.Stat(p); err == nil && !info.IsDir() {
			return p, nil
		}
	}
	return "", fmt.Errorf("not found in PATH: %s", name)
}

// WriteManifest writes a plugin manifest to the plugin directory.
func WriteManifest(pluginDir string, m PluginManifest) error {
	if m.InstalledAt == "" {
		m.InstalledAt = time.Now().UTC().Format(time.RFC3339)
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling manifest: %w", err)
	}
	return os.WriteFile(filepath.Join(pluginDir, "plugin.json"), append(data, '\n'), 0644)
}

package plugin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Install downloads and installs a plugin from a Go module path.
// It runs `go install <source>@latest`, finds the binary in GOPATH/bin,
// calls `<binary> meta` to get the manifest, and copies everything to
// ~/.human/plugins/<name>/.
func Install(source string) error {
	// Run go install.
	cmd := exec.Command("go", "install", source+"@latest")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go install failed: %s", strings.TrimSpace(stderr.String()))
	}

	// Find the binary in GOPATH/bin.
	binName := filepath.Base(source)
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("could not determine GOPATH: %w", err)
		}
		gopath = filepath.Join(home, "go")
	}
	binPath := filepath.Join(gopath, "bin", binName)
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		return fmt.Errorf("binary not found after install: %s", binPath)
	}

	return installFromPath(binPath, source)
}

// InstallFromBinary installs a plugin from a pre-built binary path.
// It calls `<binary> meta` to get the manifest, then copies the binary
// to ~/.human/plugins/<name>/.
func InstallFromBinary(binaryPath string) error {
	absPath, err := filepath.Abs(binaryPath)
	if err != nil {
		return fmt.Errorf("resolving binary path: %w", err)
	}
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("binary not found: %s", absPath)
	}
	return installFromPath(absPath, "")
}

// installFromPath is the shared logic for Install and InstallFromBinary.
func installFromPath(binPath, source string) error {
	// Get manifest from the binary.
	manifest, err := queryMeta(binPath)
	if err != nil {
		return fmt.Errorf("querying plugin metadata: %w", err)
	}
	manifest.Source = source

	// Create plugin directory.
	pluginDir := filepath.Join(PluginsDir(), manifest.Name)
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		return fmt.Errorf("creating plugin directory: %w", err)
	}

	// Copy binary to plugin directory.
	destBin := filepath.Join(pluginDir, manifest.Binary)
	data, err := os.ReadFile(binPath)
	if err != nil {
		return fmt.Errorf("reading binary: %w", err)
	}
	if err := os.WriteFile(destBin, data, 0755); err != nil {
		return fmt.Errorf("writing binary: %w", err)
	}

	// Write manifest.
	if err := WriteManifest(pluginDir, manifest); err != nil {
		return fmt.Errorf("writing manifest: %w", err)
	}

	return nil
}

// queryMeta runs `<binary> meta` and parses the JSON output into a manifest.
func queryMeta(binPath string) (PluginManifest, error) {
	cmd := exec.Command(binPath, "meta")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		if errMsg == "" {
			errMsg = err.Error()
		}
		return PluginManifest{}, fmt.Errorf("binary meta command failed: %s", errMsg)
	}

	var m PluginManifest
	if err := json.Unmarshal(stdout.Bytes(), &m); err != nil {
		return PluginManifest{}, fmt.Errorf("parsing meta output: %w", err)
	}

	if m.Name == "" {
		return PluginManifest{}, fmt.Errorf("plugin meta missing name")
	}

	// Set the binary name from the file name.
	m.Binary = filepath.Base(binPath)

	return m, nil
}

// Uninstall removes an installed plugin by deleting its directory.
func Uninstall(name string) error {
	pluginDir := filepath.Join(PluginsDir(), name)
	if _, err := os.Stat(pluginDir); os.IsNotExist(err) {
		return fmt.Errorf("plugin not found: %s", name)
	}
	return os.RemoveAll(pluginDir)
}

// List returns all installed plugin manifests. This is an alias for Discover.
func List() ([]PluginManifest, error) {
	return Discover()
}

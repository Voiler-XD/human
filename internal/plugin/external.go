package plugin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/barun-bash/human/internal/codegen"
	"github.com/barun-bash/human/internal/ir"
)

// ExternalGenerator adapts an external plugin binary to the CodeGenerator
// interface. It communicates with the plugin via CLI subcommands and JSON.
type ExternalGenerator struct {
	manifest PluginManifest
	binary   string
	settings map[string]string
}

// Meta returns the generator's metadata derived from its manifest.
func (g *ExternalGenerator) Meta() codegen.PluginMeta {
	return codegen.PluginMeta{
		Name:        g.manifest.Name,
		Version:     g.manifest.Version,
		Description: g.manifest.Description,
		Category:    codegen.Category(g.manifest.Category),
	}
}

// Enabled always returns true for external plugins. Disabling is handled
// by the config system's tri-state override.
func (g *ExternalGenerator) Enabled(_ *ir.Application) bool {
	return true
}

// StageName returns the display name for progress reporting.
func (g *ExternalGenerator) StageName() string {
	return fmt.Sprintf("Running plugin: %s", g.manifest.Name)
}

// OutputDir returns the subdirectory name for this plugin's output.
func (g *ExternalGenerator) OutputDir() string {
	return g.manifest.Name
}

// SetSettings injects configuration settings from the project config.
func (g *ExternalGenerator) SetSettings(settings map[string]string) {
	g.settings = settings
}

// Generate writes the IR to a temp file, invokes the plugin binary with
// the generate subcommand, and captures any errors from stderr.
func (g *ExternalGenerator) Generate(app *ir.Application, outputDir string) error {
	// Write IR to a temporary file.
	irData, err := json.Marshal(app)
	if err != nil {
		return fmt.Errorf("marshaling IR for plugin %s: %w", g.manifest.Name, err)
	}

	tmpFile, err := os.CreateTemp("", "human-ir-*.json")
	if err != nil {
		return fmt.Errorf("creating temp IR file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(irData); err != nil {
		tmpFile.Close()
		return fmt.Errorf("writing temp IR file: %w", err)
	}
	tmpFile.Close()

	// Resolve output directory to an absolute path so the plugin writes
	// to the correct location regardless of its working directory.
	absOutput, err := filepath.Abs(outputDir)
	if err != nil {
		return fmt.Errorf("resolving output dir: %w", err)
	}

	// Ensure output directory exists.
	if err := os.MkdirAll(absOutput, 0755); err != nil {
		return fmt.Errorf("creating output dir for plugin %s: %w", g.manifest.Name, err)
	}

	// Build command arguments.
	args := []string{"generate", "--ir", tmpFile.Name(), "--output", absOutput}
	if len(g.settings) > 0 {
		settingsJSON, err := json.Marshal(g.settings)
		if err != nil {
			return fmt.Errorf("marshaling settings for plugin %s: %w", g.manifest.Name, err)
		}
		args = append(args, "--settings", string(settingsJSON))
	}

	// Execute the plugin binary.
	cmd := exec.Command(g.binary, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errMsg := stderr.String()
		if errMsg == "" {
			errMsg = err.Error()
		}
		return fmt.Errorf("plugin %s failed: %s", g.manifest.Name, errMsg)
	}

	return nil
}

// BinaryPath returns the resolved path to the plugin binary.
// Exported for testing.
func (g *ExternalGenerator) BinaryPath() string {
	return g.binary
}

// NewExternalGenerator creates an ExternalGenerator from a manifest and
// resolved binary path. Exported for testing.
func NewExternalGenerator(manifest PluginManifest, binaryPath string) *ExternalGenerator {
	return &ExternalGenerator{
		manifest: manifest,
		binary:   binaryPath,
	}
}

// resolveOutputDir returns the full output path for a plugin.
func resolveOutputDir(baseDir string, pluginName string) string {
	return filepath.Join(baseDir, pluginName)
}

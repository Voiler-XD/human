package plugin

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPluginsDir(t *testing.T) {
	dir := PluginsDir()
	if dir == "" {
		t.Fatal("PluginsDir returned empty string")
	}
	if !strings.Contains(dir, "plugins") {
		t.Errorf("PluginsDir should contain 'plugins', got %q", dir)
	}
}

func TestDiscoverEmpty(t *testing.T) {
	// Override PluginsDir by using a temp dir that doesn't exist.
	orig := os.Getenv("HOME")
	tmp := t.TempDir()
	os.Setenv("HOME", tmp)
	defer os.Setenv("HOME", orig)

	manifests, err := Discover()
	if err != nil {
		t.Fatalf("Discover on empty dir: %v", err)
	}
	if len(manifests) != 0 {
		t.Errorf("expected 0 manifests, got %d", len(manifests))
	}
}

func TestDiscoverWithPlugins(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("HOME", tmp)
	defer os.Unsetenv("HOME")

	// Create a valid plugin.
	pluginDir := filepath.Join(tmp, ".human", "plugins", "test-plugin")
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		t.Fatal(err)
	}
	manifest := `{"name":"test-plugin","version":"1.0.0","description":"Test","category":"frontend","binary":"test-plugin"}`
	if err := os.WriteFile(filepath.Join(pluginDir, "plugin.json"), []byte(manifest), 0644); err != nil {
		t.Fatal(err)
	}

	manifests, err := Discover()
	if err != nil {
		t.Fatalf("Discover: %v", err)
	}
	if len(manifests) != 1 {
		t.Fatalf("expected 1 manifest, got %d", len(manifests))
	}
	if manifests[0].Name != "test-plugin" {
		t.Errorf("expected name 'test-plugin', got %q", manifests[0].Name)
	}
}

func TestDiscoverSkipsBroken(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("HOME", tmp)
	defer os.Unsetenv("HOME")

	pluginsDir := filepath.Join(tmp, ".human", "plugins")

	// Create a valid plugin.
	goodDir := filepath.Join(pluginsDir, "good-plugin")
	if err := os.MkdirAll(goodDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(goodDir, "plugin.json"),
		[]byte(`{"name":"good-plugin","version":"1.0.0","description":"Good","category":"backend","binary":"good"}`), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a broken plugin (invalid JSON).
	brokenDir := filepath.Join(pluginsDir, "broken-plugin")
	if err := os.MkdirAll(brokenDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(brokenDir, "plugin.json"), []byte(`{invalid`), 0644); err != nil {
		t.Fatal(err)
	}

	manifests, err := Discover()
	if err != nil {
		t.Fatalf("Discover: %v", err)
	}
	if len(manifests) != 1 {
		t.Fatalf("expected 1 manifest (skipping broken), got %d", len(manifests))
	}
	if manifests[0].Name != "good-plugin" {
		t.Errorf("expected name 'good-plugin', got %q", manifests[0].Name)
	}
}

func TestLoadManifestValid(t *testing.T) {
	tmp := t.TempDir()
	manifest := `{"name":"my-plugin","version":"2.0.0","description":"My Plugin","category":"infra","binary":"my-plugin"}`
	if err := os.WriteFile(filepath.Join(tmp, "plugin.json"), []byte(manifest), 0644); err != nil {
		t.Fatal(err)
	}

	m, err := LoadManifest(tmp)
	if err != nil {
		t.Fatalf("LoadManifest: %v", err)
	}
	if m.Name != "my-plugin" {
		t.Errorf("expected name 'my-plugin', got %q", m.Name)
	}
	if m.Version != "2.0.0" {
		t.Errorf("expected version '2.0.0', got %q", m.Version)
	}
	if m.Category != "infra" {
		t.Errorf("expected category 'infra', got %q", m.Category)
	}
}

func TestLoadManifestInvalid(t *testing.T) {
	tmp := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmp, "plugin.json"), []byte(`not json`), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadManifest(tmp)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestLoadManifestMissingName(t *testing.T) {
	tmp := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmp, "plugin.json"), []byte(`{"version":"1.0.0"}`), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadManifest(tmp)
	if err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestLoadManifestNoFile(t *testing.T) {
	tmp := t.TempDir()
	_, err := LoadManifest(tmp)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestResolveBinaryLocal(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("HOME", tmp)
	defer os.Unsetenv("HOME")

	pluginDir := filepath.Join(tmp, ".human", "plugins", "local-bin")
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		t.Fatal(err)
	}
	// Create a fake binary.
	binPath := filepath.Join(pluginDir, "local-bin")
	if err := os.WriteFile(binPath, []byte("#!/bin/bash\n"), 0755); err != nil {
		t.Fatal(err)
	}

	m := PluginManifest{Name: "local-bin", Binary: "local-bin"}
	result := resolveBinary(m)
	if result != binPath {
		t.Errorf("expected %q, got %q", binPath, result)
	}
}

func TestResolveBinaryNotFound(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("HOME", tmp)
	defer os.Unsetenv("HOME")

	m := PluginManifest{Name: "nonexistent", Binary: "nonexistent-binary-xyz"}
	result := resolveBinary(m)
	if result != "" {
		t.Errorf("expected empty string for missing binary, got %q", result)
	}
}

func TestWriteManifest(t *testing.T) {
	tmp := t.TempDir()
	m := PluginManifest{
		Name:        "write-test",
		Version:     "1.0.0",
		Description: "Write test",
		Category:    "backend",
		Binary:      "write-test",
	}

	if err := WriteManifest(tmp, m); err != nil {
		t.Fatalf("WriteManifest: %v", err)
	}

	// Read it back.
	loaded, err := LoadManifest(tmp)
	if err != nil {
		t.Fatalf("LoadManifest after write: %v", err)
	}
	if loaded.Name != m.Name {
		t.Errorf("expected name %q, got %q", m.Name, loaded.Name)
	}
	if loaded.InstalledAt == "" {
		t.Error("expected InstalledAt to be set")
	}
}

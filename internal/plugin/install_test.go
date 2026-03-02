package plugin

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestInstallFromBinary(t *testing.T) {
	bin := mockPluginBinary(t)
	if _, err := os.Stat(bin); os.IsNotExist(err) {
		t.Skip("mock-plugin.sh not found")
	}

	tmp := t.TempDir()
	os.Setenv("HOME", tmp)
	defer os.Unsetenv("HOME")

	if err := InstallFromBinary(bin); err != nil {
		t.Fatalf("InstallFromBinary: %v", err)
	}

	// Verify plugin was installed.
	pluginDir := filepath.Join(tmp, ".human", "plugins", "mock-plugin")
	if _, err := os.Stat(filepath.Join(pluginDir, "plugin.json")); os.IsNotExist(err) {
		t.Error("expected plugin.json in plugin directory")
	}

	// Verify the binary was copied.
	binName := "mock-plugin.sh"
	if _, err := os.Stat(filepath.Join(pluginDir, binName)); os.IsNotExist(err) {
		t.Errorf("expected binary %s in plugin directory", binName)
	}

	// Verify manifest contents.
	m, err := LoadManifest(pluginDir)
	if err != nil {
		t.Fatalf("LoadManifest: %v", err)
	}
	if m.Name != "mock-plugin" {
		t.Errorf("expected name 'mock-plugin', got %q", m.Name)
	}
	if m.Version != "0.1.0" {
		t.Errorf("expected version '0.1.0', got %q", m.Version)
	}
}

func TestInstallFromBinaryNotFound(t *testing.T) {
	err := InstallFromBinary("/nonexistent/path/to/binary")
	if err == nil {
		t.Fatal("expected error for nonexistent binary")
	}
}

func TestInstallInvalidBinary(t *testing.T) {
	// Create a binary that doesn't support the 'meta' command.
	tmp := t.TempDir()
	os.Setenv("HOME", tmp)
	defer os.Unsetenv("HOME")

	fakeBin := filepath.Join(tmp, "fake-plugin")
	if runtime.GOOS == "windows" {
		fakeBin += ".exe"
	}
	os.WriteFile(fakeBin, []byte("#!/bin/bash\nexit 1\n"), 0755)

	err := InstallFromBinary(fakeBin)
	if err == nil {
		t.Fatal("expected error for invalid binary")
	}
}

func TestUninstall(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("HOME", tmp)
	defer os.Unsetenv("HOME")

	// Create a fake plugin.
	pluginDir := filepath.Join(tmp, ".human", "plugins", "removeme")
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(pluginDir, "plugin.json"),
		[]byte(`{"name":"removeme","version":"1.0.0"}`), 0644)

	if err := Uninstall("removeme"); err != nil {
		t.Fatalf("Uninstall: %v", err)
	}

	if _, err := os.Stat(pluginDir); !os.IsNotExist(err) {
		t.Error("expected plugin directory to be removed")
	}
}

func TestUninstallNotFound(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("HOME", tmp)
	defer os.Unsetenv("HOME")

	err := Uninstall("nonexistent-plugin")
	if err == nil {
		t.Fatal("expected error for nonexistent plugin")
	}
}

func TestWriteAndLoadManifestRoundTrip(t *testing.T) {
	tmp := t.TempDir()
	m := PluginManifest{
		Name:        "roundtrip",
		Version:     "2.3.4",
		Description: "Round trip test",
		Category:    "infra",
		Binary:      "roundtrip-bin",
		Source:      "github.com/example/roundtrip",
	}

	if err := WriteManifest(tmp, m); err != nil {
		t.Fatalf("WriteManifest: %v", err)
	}

	loaded, err := LoadManifest(tmp)
	if err != nil {
		t.Fatalf("LoadManifest: %v", err)
	}

	if loaded.Name != m.Name || loaded.Version != m.Version || loaded.Source != m.Source {
		t.Errorf("round-trip mismatch: got %+v", loaded)
	}
	if loaded.InstalledAt == "" {
		t.Error("InstalledAt should be populated")
	}
}

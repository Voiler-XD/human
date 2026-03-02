package plugin

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/barun-bash/human/internal/codegen"
	"github.com/barun-bash/human/internal/ir"
)

func mockPluginBinary(t *testing.T) string {
	t.Helper()
	// Get the path to the mock plugin script.
	_, thisFile, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(thisFile), "testdata", "mock-plugin.sh")
}

func TestExternalGeneratorMeta(t *testing.T) {
	g := NewExternalGenerator(PluginManifest{
		Name:        "test-gen",
		Version:     "1.2.3",
		Description: "Test generator",
		Category:    "frontend",
	}, "/fake/path")

	meta := g.Meta()
	if meta.Name != "test-gen" {
		t.Errorf("expected name 'test-gen', got %q", meta.Name)
	}
	if meta.Version != "1.2.3" {
		t.Errorf("expected version '1.2.3', got %q", meta.Version)
	}
	if meta.Category != codegen.CategoryFrontend {
		t.Errorf("expected category 'frontend', got %q", meta.Category)
	}
}

func TestExternalGeneratorEnabled(t *testing.T) {
	g := NewExternalGenerator(PluginManifest{Name: "x"}, "/fake/path")
	if !g.Enabled(nil) {
		t.Error("ExternalGenerator should always be enabled")
	}
}

func TestExternalGeneratorStageName(t *testing.T) {
	g := NewExternalGenerator(PluginManifest{Name: "htmx"}, "/fake/path")
	expected := "Running plugin: htmx"
	if got := g.StageName(); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestExternalGeneratorOutputDir(t *testing.T) {
	g := NewExternalGenerator(PluginManifest{Name: "k8s"}, "/fake/path")
	if got := g.OutputDir(); got != "k8s" {
		t.Errorf("expected 'k8s', got %q", got)
	}
}

func TestExternalGeneratorGenerate(t *testing.T) {
	bin := mockPluginBinary(t)
	if _, err := os.Stat(bin); os.IsNotExist(err) {
		t.Skip("mock-plugin.sh not found")
	}

	app := &ir.Application{
		Name:     "TestApp",
		Platform: "web",
		Config: &ir.BuildConfig{
			Frontend: "HTMX",
		},
		Pages: []*ir.Page{{
			Name: "Home",
		}},
	}

	outputDir := t.TempDir()
	g := NewExternalGenerator(PluginManifest{
		Name:    "mock-plugin",
		Version: "0.1.0",
		Binary:  "mock-plugin.sh",
	}, bin)

	if err := g.Generate(app, outputDir); err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Verify output file was created.
	indexPath := filepath.Join(outputDir, "index.html")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		t.Error("expected index.html to be created")
	}
	content, _ := os.ReadFile(indexPath)
	if !strings.Contains(string(content), "mock-plugin") {
		t.Error("index.html should contain mock-plugin marker")
	}
}

func TestExternalGeneratorGenerateFailure(t *testing.T) {
	bin := mockPluginBinary(t)
	if _, err := os.Stat(bin); os.IsNotExist(err) {
		t.Skip("mock-plugin.sh not found")
	}

	// Create a generator with the mock binary but pass an unknown command
	// by having an invalid state. We'll test by using a binary that fails.
	failScript := filepath.Join(t.TempDir(), "fail.sh")
	os.WriteFile(failScript, []byte("#!/bin/bash\necho 'plugin error' >&2\nexit 1\n"), 0755)

	g := NewExternalGenerator(PluginManifest{Name: "fail-plugin"}, failScript)
	err := g.Generate(&ir.Application{Name: "Test"}, t.TempDir())
	if err == nil {
		t.Fatal("expected error from failing plugin")
	}
	if !strings.Contains(err.Error(), "fail-plugin") {
		t.Errorf("error should mention plugin name, got: %v", err)
	}
}

func TestExternalGeneratorIRRoundTrip(t *testing.T) {
	bin := mockPluginBinary(t)
	if _, err := os.Stat(bin); os.IsNotExist(err) {
		t.Skip("mock-plugin.sh not found")
	}

	app := &ir.Application{
		Name:     "RoundTripApp",
		Platform: "web",
		Config: &ir.BuildConfig{
			Frontend: "React",
			Backend:  "Node",
			Database: "PostgreSQL",
		},
		Data: []*ir.DataModel{{
			Name: "User",
			Fields: []*ir.DataField{
				{Name: "name", Type: "text", Required: true},
				{Name: "email", Type: "email", Required: true, Unique: true},
			},
		}},
	}

	outputDir := t.TempDir()
	g := NewExternalGenerator(PluginManifest{
		Name:    "mock-plugin",
		Version: "0.1.0",
		Binary:  "mock-plugin.sh",
	}, bin)

	if err := g.Generate(app, outputDir); err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Read the IR copy the mock plugin wrote.
	irCopy, err := os.ReadFile(filepath.Join(outputDir, "ir-copy.json"))
	if err != nil {
		t.Fatalf("reading IR copy: %v", err)
	}

	var decoded ir.Application
	if err := json.Unmarshal(irCopy, &decoded); err != nil {
		t.Fatalf("decoding IR copy: %v", err)
	}

	if decoded.Name != "RoundTripApp" {
		t.Errorf("expected name 'RoundTripApp', got %q", decoded.Name)
	}
	if len(decoded.Data) != 1 || decoded.Data[0].Name != "User" {
		t.Error("IR round-trip lost data models")
	}
}

func TestExternalGeneratorSettings(t *testing.T) {
	bin := mockPluginBinary(t)
	if _, err := os.Stat(bin); os.IsNotExist(err) {
		t.Skip("mock-plugin.sh not found")
	}

	outputDir := t.TempDir()
	g := NewExternalGenerator(PluginManifest{
		Name:    "mock-plugin",
		Version: "0.1.0",
		Binary:  "mock-plugin.sh",
	}, bin)

	g.SetSettings(map[string]string{"key": "value", "another": "setting"})

	if err := g.Generate(&ir.Application{Name: "SettingsTest"}, outputDir); err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Verify settings were passed through.
	settingsData, err := os.ReadFile(filepath.Join(outputDir, "settings.json"))
	if err != nil {
		t.Fatalf("reading settings.json: %v", err)
	}

	var settings map[string]string
	if err := json.Unmarshal(settingsData, &settings); err != nil {
		t.Fatalf("decoding settings: %v", err)
	}
	if settings["key"] != "value" {
		t.Errorf("expected key=value, got key=%q", settings["key"])
	}
}

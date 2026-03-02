package plugin

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestScaffoldCreatesExpectedFiles(t *testing.T) {
	tmp := t.TempDir()
	outputDir := filepath.Join(tmp, "test-plugin")

	if err := Scaffold("test-plugin", "frontend", outputDir); err != nil {
		t.Fatalf("Scaffold: %v", err)
	}

	expected := []string{"main.go", "generator.go", "go.mod", "Makefile", "README.md"}
	for _, f := range expected {
		path := filepath.Join(outputDir, f)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected file %s not found", f)
		}
	}
}

func TestScaffoldMainContainsMetaAndGenerate(t *testing.T) {
	tmp := t.TempDir()
	outputDir := filepath.Join(tmp, "my-gen")

	if err := Scaffold("my-gen", "backend", outputDir); err != nil {
		t.Fatalf("Scaffold: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(outputDir, "main.go"))
	if err != nil {
		t.Fatalf("reading main.go: %v", err)
	}

	s := string(content)
	if !strings.Contains(s, `"meta"`) {
		t.Error("main.go should contain 'meta' case")
	}
	if !strings.Contains(s, `"generate"`) {
		t.Error("main.go should contain 'generate' case")
	}
	if !strings.Contains(s, `"my-gen"`) {
		t.Error("main.go should contain plugin name")
	}
}

func TestScaffoldGoModCorrect(t *testing.T) {
	tmp := t.TempDir()
	outputDir := filepath.Join(tmp, "k8s-gen")

	if err := Scaffold("k8s-gen", "infra", outputDir); err != nil {
		t.Fatalf("Scaffold: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(outputDir, "go.mod"))
	if err != nil {
		t.Fatalf("reading go.mod: %v", err)
	}

	s := string(content)
	if !strings.Contains(s, "human-plugin-k8s-gen") {
		t.Error("go.mod should contain plugin module name")
	}
	if !strings.Contains(s, "go 1.25") {
		t.Error("go.mod should specify go version")
	}
}

func TestScaffoldCustomNamePropagated(t *testing.T) {
	tmp := t.TempDir()
	outputDir := filepath.Join(tmp, "custom-name")

	if err := Scaffold("custom-name", "database", outputDir); err != nil {
		t.Fatalf("Scaffold: %v", err)
	}

	// Check that the name appears in the README.
	readme, err := os.ReadFile(filepath.Join(outputDir, "README.md"))
	if err != nil {
		t.Fatalf("reading README.md: %v", err)
	}
	if !strings.Contains(string(readme), "custom-name") {
		t.Error("README should contain the plugin name")
	}

	// Check the Makefile references the plugin name.
	makefile, err := os.ReadFile(filepath.Join(outputDir, "Makefile"))
	if err != nil {
		t.Fatalf("reading Makefile: %v", err)
	}
	if !strings.Contains(string(makefile), "custom-name") {
		t.Error("Makefile should contain the plugin name")
	}
}

func TestScaffoldInvalidCategory(t *testing.T) {
	tmp := t.TempDir()
	err := Scaffold("bad-cat", "invalid", filepath.Join(tmp, "bad"))
	if err == nil {
		t.Fatal("expected error for invalid category")
	}
}

func TestScaffoldEmptyName(t *testing.T) {
	err := Scaffold("", "frontend", t.TempDir())
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

package parser_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/barun-bash/human/internal/ir"
	"github.com/barun-bash/human/internal/parser"
)

// projectRoot returns the repository root by walking up from the test file.
func projectRoot(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot determine test file path")
	}
	// filename is .../internal/parser/merge_integration_test.go
	// root is two directories up
	return filepath.Join(filepath.Dir(filename), "..", "..")
}

// TestMergeProducesIdenticalIR is the gold standard integration test.
// It parses the single-file taskflow example and the multi-file taskflow-multi
// example, builds IR from each, serializes to YAML, and compares.
func TestMergeProducesIdenticalIR(t *testing.T) {
	root := projectRoot(t)

	singleFile := filepath.Join(root, "examples", "taskflow", "app.human")
	multiDir := filepath.Join(root, "examples", "taskflow-multi")

	// Verify both exist.
	if _, err := os.Stat(singleFile); err != nil {
		t.Skipf("single-file example not found: %v", err)
	}
	if _, err := os.Stat(multiDir); err != nil {
		t.Skipf("multi-file example not found: %v", err)
	}

	// Parse single-file.
	singleSource, err := os.ReadFile(singleFile)
	if err != nil {
		t.Fatalf("reading single file: %v", err)
	}
	singleProg, err := parser.Parse(string(singleSource))
	if err != nil {
		t.Fatalf("parsing single file: %v", err)
	}

	// Parse multi-file.
	multiFiles, err := parser.DiscoverFiles(multiDir)
	if err != nil {
		t.Fatalf("discovering multi files: %v", err)
	}
	if len(multiFiles) < 2 {
		t.Fatalf("expected multiple files, got %d", len(multiFiles))
	}

	multiPrograms, err := parser.ParseFiles(multiFiles)
	if err != nil {
		t.Fatalf("parsing multi files: %v", err)
	}

	multiProg, err := parser.MergePrograms(multiPrograms)
	if err != nil {
		t.Fatalf("merging programs: %v", err)
	}

	// Build IR from both.
	singleApp, err := ir.Build(singleProg)
	if err != nil {
		t.Fatalf("building single IR: %v", err)
	}
	multiApp, err := ir.Build(multiProg)
	if err != nil {
		t.Fatalf("building multi IR: %v", err)
	}

	// Serialize both to YAML and compare.
	singleYAML, err := ir.ToYAML(singleApp)
	if err != nil {
		t.Fatalf("serializing single IR: %v", err)
	}
	multiYAML, err := ir.ToYAML(multiApp)
	if err != nil {
		t.Fatalf("serializing multi IR: %v", err)
	}

	if singleYAML != multiYAML {
		// Find the first difference for a useful error message.
		sLines := splitLines(singleYAML)
		mLines := splitLines(multiYAML)
		for i := 0; i < len(sLines) && i < len(mLines); i++ {
			if sLines[i] != mLines[i] {
				t.Fatalf("IR mismatch at line %d:\n  single: %s\n  multi:  %s", i+1, sLines[i], mLines[i])
			}
		}
		if len(sLines) != len(mLines) {
			t.Fatalf("IR line count mismatch: single=%d, multi=%d", len(sLines), len(mLines))
		}
		t.Fatal("IR YAML content differs (unknown line)")
	}
}

// TestMultiFileDiscovery verifies that the multi-file example is properly discovered.
func TestMultiFileDiscovery(t *testing.T) {
	root := projectRoot(t)
	multiDir := filepath.Join(root, "examples", "taskflow-multi")

	if _, err := os.Stat(multiDir); err != nil {
		t.Skipf("multi-file example not found: %v", err)
	}

	files, err := parser.DiscoverFiles(multiDir)
	if err != nil {
		t.Fatal(err)
	}

	if len(files) < 2 {
		t.Fatalf("expected multiple files, got %d", len(files))
	}

	// app.human must be first
	if filepath.Base(files[0]) != "app.human" {
		t.Errorf("expected app.human first, got %s", filepath.Base(files[0]))
	}
}

// TestMultiFileParseAndMerge verifies that all declarations from multi-file
// are properly merged into a single program.
func TestMultiFileParseAndMerge(t *testing.T) {
	root := projectRoot(t)
	multiDir := filepath.Join(root, "examples", "taskflow-multi")

	if _, err := os.Stat(multiDir); err != nil {
		t.Skipf("multi-file example not found: %v", err)
	}

	files, err := parser.DiscoverFiles(multiDir)
	if err != nil {
		t.Fatal(err)
	}

	programs, err := parser.ParseFiles(files)
	if err != nil {
		t.Fatal(err)
	}

	merged, err := parser.MergePrograms(programs)
	if err != nil {
		t.Fatal(err)
	}

	// Verify key declarations were merged.
	if merged.App == nil {
		t.Error("expected app declaration")
	}
	if len(merged.Data) == 0 {
		t.Error("expected data declarations")
	}
	if len(merged.Pages) == 0 {
		t.Error("expected page declarations")
	}
	if len(merged.APIs) == 0 {
		t.Error("expected API declarations")
	}
	if merged.Build == nil {
		t.Error("expected build declaration")
	}
	if merged.Theme == nil {
		t.Error("expected theme declaration")
	}
	if merged.Authentication == nil {
		t.Error("expected authentication declaration")
	}
	if merged.Database == nil {
		t.Error("expected database declaration")
	}
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

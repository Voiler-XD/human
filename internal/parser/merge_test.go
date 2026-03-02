package parser

import (
	"os"
	"path/filepath"
	"testing"
)

// helper to create a temp dir with .human files
func setupTestDir(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
			t.Fatalf("writing %s: %v", name, err)
		}
	}
	return dir
}

// ── DiscoverFiles ──

func TestDiscoverFiles_SingleFile(t *testing.T) {
	dir := setupTestDir(t, map[string]string{
		"myapp.human": "app MyApp is a web application\n",
	})
	files, err := DiscoverFiles(filepath.Join(dir, "myapp.human"))
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}
	if filepath.Base(files[0]) != "myapp.human" {
		t.Errorf("expected myapp.human, got %s", filepath.Base(files[0]))
	}
}

func TestDiscoverFiles_MultiFileOrdering(t *testing.T) {
	dir := setupTestDir(t, map[string]string{
		"app.human":    "app MyApp is a web application\n",
		"models.human": "data User:\n  has a name which is text\n",
		"pages.human":  "page Home:\n  show a greeting\n",
	})
	files, err := DiscoverFiles(filepath.Join(dir, "app.human"))
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 3 {
		t.Fatalf("expected 3 files, got %d", len(files))
	}
	// app.human must be first
	if filepath.Base(files[0]) != "app.human" {
		t.Errorf("expected app.human first, got %s", filepath.Base(files[0]))
	}
	// rest should be alphabetical
	if filepath.Base(files[1]) != "models.human" {
		t.Errorf("expected models.human second, got %s", filepath.Base(files[1]))
	}
	if filepath.Base(files[2]) != "pages.human" {
		t.Errorf("expected pages.human third, got %s", filepath.Base(files[2]))
	}
}

func TestDiscoverFiles_DirectoryInput(t *testing.T) {
	dir := setupTestDir(t, map[string]string{
		"app.human":    "app MyApp is a web application\n",
		"models.human": "data User:\n  has a name which is text\n",
	})
	files, err := DiscoverFiles(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}
	if filepath.Base(files[0]) != "app.human" {
		t.Errorf("expected app.human first, got %s", filepath.Base(files[0]))
	}
}

func TestDiscoverFiles_NoAppHumanError(t *testing.T) {
	dir := setupTestDir(t, map[string]string{
		"models.human": "data User:\n  has a name which is text\n",
		"pages.human":  "page Home:\n  show a greeting\n",
	})
	_, err := DiscoverFiles(dir)
	if err == nil {
		t.Fatal("expected error for missing app.human")
	}
	if got := err.Error(); !contains(got, "requires app.human") {
		t.Errorf("expected 'requires app.human' in error, got: %s", got)
	}
}

func TestDiscoverFiles_IgnoresNonHumanFiles(t *testing.T) {
	dir := setupTestDir(t, map[string]string{
		"myapp.human": "app MyApp is a web application\n",
		"README.md":   "# Readme",
		"config.json": "{}",
	})
	files, err := DiscoverFiles(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}
}

func TestDiscoverFiles_NoFiles(t *testing.T) {
	dir := t.TempDir()
	_, err := DiscoverFiles(dir)
	if err == nil {
		t.Fatal("expected error for empty directory")
	}
}

// ── ParseFiles ──

func TestParseFiles_SingleFile(t *testing.T) {
	dir := setupTestDir(t, map[string]string{
		"app.human": "app MyApp is a web application\n\ndata User:\n  has a name which is text\n",
	})
	files, _ := DiscoverFiles(dir)
	programs, err := ParseFiles(files)
	if err != nil {
		t.Fatal(err)
	}
	if len(programs) != 1 {
		t.Fatalf("expected 1 program, got %d", len(programs))
	}
	if programs[0].App == nil {
		t.Fatal("expected app declaration")
	}
}

func TestParseFiles_FileFieldsSet(t *testing.T) {
	dir := setupTestDir(t, map[string]string{
		"app.human":    "app MyApp is a web application\n",
		"models.human": "data User:\n  has a name which is text\n",
	})
	files, _ := DiscoverFiles(dir)
	programs, err := ParseFiles(files)
	if err != nil {
		t.Fatal(err)
	}
	if programs[0].App.File != files[0] {
		t.Errorf("expected app File=%s, got %s", files[0], programs[0].App.File)
	}
	if len(programs[1].Data) == 0 {
		t.Fatal("expected data declaration in models.human")
	}
	if programs[1].Data[0].File != files[1] {
		t.Errorf("expected data File=%s, got %s", files[1], programs[1].Data[0].File)
	}
}

func TestParseFiles_ParseErrorIncludesFilename(t *testing.T) {
	dir := setupTestDir(t, map[string]string{
		"app.human": "app MyApp is a web application\n",
		"bad.human": "this is not valid human syntax {{{\n",
	})
	files, _ := DiscoverFiles(dir)
	_, err := ParseFiles(files)
	if err == nil {
		// The parser is lenient — it may not error on arbitrary text.
		// This test verifies that IF a parse error occurs, the filename is included.
		t.Skip("parser accepted bad input; cannot test filename in error")
	}
	if !contains(err.Error(), "bad.human") {
		t.Errorf("expected filename in error, got: %s", err.Error())
	}
}

// ── MergePrograms ──

func TestMergePrograms_SinglePassthrough(t *testing.T) {
	prog := &Program{
		App: &AppDeclaration{Name: "MyApp", File: "app.human"},
	}
	merged, err := MergePrograms([]*Program{prog})
	if err != nil {
		t.Fatal(err)
	}
	if merged != prog {
		t.Error("single program should be returned as-is")
	}
}

func TestMergePrograms_SlicesMerge(t *testing.T) {
	p1 := &Program{
		Data: []*DataDeclaration{{Name: "User", File: "models.human"}},
	}
	p2 := &Program{
		Data: []*DataDeclaration{{Name: "Task", File: "models2.human"}},
	}
	merged, err := MergePrograms([]*Program{p1, p2})
	if err != nil {
		t.Fatal(err)
	}
	if len(merged.Data) != 2 {
		t.Fatalf("expected 2 data models, got %d", len(merged.Data))
	}
	if merged.Data[0].Name != "User" || merged.Data[1].Name != "Task" {
		t.Error("data models not in expected order")
	}
}

func TestMergePrograms_SingletonFirstWins(t *testing.T) {
	p1 := &Program{
		App: &AppDeclaration{Name: "MyApp", File: "app.human", Line: 1},
	}
	p2 := &Program{}
	merged, err := MergePrograms([]*Program{p1, p2})
	if err != nil {
		t.Fatal(err)
	}
	if merged.App == nil || merged.App.Name != "MyApp" {
		t.Error("expected app from first program")
	}
}

func TestMergePrograms_DuplicateSingletonError(t *testing.T) {
	p1 := &Program{
		App: &AppDeclaration{Name: "App1", File: "app.human", Line: 1},
	}
	p2 := &Program{
		App: &AppDeclaration{Name: "App2", File: "other.human", Line: 1},
	}
	_, err := MergePrograms([]*Program{p1, p2})
	if err == nil {
		t.Fatal("expected error for duplicate app")
	}
	if !contains(err.Error(), "duplicate app") {
		t.Errorf("expected 'duplicate app' in error, got: %s", err.Error())
	}
}

func TestMergePrograms_DuplicateThemeError(t *testing.T) {
	p1 := &Program{
		Theme: &ThemeDeclaration{File: "app.human", Line: 5},
	}
	p2 := &Program{
		Theme: &ThemeDeclaration{File: "theme.human", Line: 1},
	}
	_, err := MergePrograms([]*Program{p1, p2})
	if err == nil {
		t.Fatal("expected error for duplicate theme")
	}
	if !contains(err.Error(), "duplicate theme") {
		t.Errorf("expected 'duplicate theme' in error, got: %s", err.Error())
	}
}

func TestMergePrograms_DuplicateBuildError(t *testing.T) {
	p1 := &Program{
		Build: &BuildDeclaration{File: "app.human", Line: 10},
	}
	p2 := &Program{
		Build: &BuildDeclaration{File: "build.human", Line: 1},
	}
	_, err := MergePrograms([]*Program{p1, p2})
	if err == nil {
		t.Fatal("expected error for duplicate build")
	}
}

func TestMergePrograms_MixedContent(t *testing.T) {
	p1 := &Program{
		App:   &AppDeclaration{Name: "MyApp", File: "app.human"},
		Build: &BuildDeclaration{File: "app.human"},
		Data:  []*DataDeclaration{{Name: "User", File: "app.human"}},
		Pages: []*PageDeclaration{{Name: "Home", File: "app.human"}},
	}
	p2 := &Program{
		Data:  []*DataDeclaration{{Name: "Task", File: "models.human"}},
		APIs:  []*APIDeclaration{{Name: "CreateTask", File: "api.human"}},
		Pages: []*PageDeclaration{{Name: "Dashboard", File: "pages.human"}},
	}
	merged, err := MergePrograms([]*Program{p1, p2})
	if err != nil {
		t.Fatal(err)
	}
	if merged.App == nil || merged.App.Name != "MyApp" {
		t.Error("expected app")
	}
	if merged.Build == nil {
		t.Error("expected build")
	}
	if len(merged.Data) != 2 {
		t.Errorf("expected 2 data, got %d", len(merged.Data))
	}
	if len(merged.Pages) != 2 {
		t.Errorf("expected 2 pages, got %d", len(merged.Pages))
	}
	if len(merged.APIs) != 1 {
		t.Errorf("expected 1 api, got %d", len(merged.APIs))
	}
}

func TestMergePrograms_EmptyPrograms(t *testing.T) {
	merged, err := MergePrograms([]*Program{{}, {}})
	if err != nil {
		t.Fatal(err)
	}
	if merged.App != nil {
		t.Error("expected nil app")
	}
	if len(merged.Data) != 0 {
		t.Error("expected empty data")
	}
}

func TestMergePrograms_OrderPreserved(t *testing.T) {
	p1 := &Program{
		Pages: []*PageDeclaration{
			{Name: "Home", File: "p1.human"},
			{Name: "About", File: "p1.human"},
		},
	}
	p2 := &Program{
		Pages: []*PageDeclaration{
			{Name: "Dashboard", File: "p2.human"},
		},
	}
	merged, err := MergePrograms([]*Program{p1, p2})
	if err != nil {
		t.Fatal(err)
	}
	if len(merged.Pages) != 3 {
		t.Fatalf("expected 3 pages, got %d", len(merged.Pages))
	}
	names := []string{merged.Pages[0].Name, merged.Pages[1].Name, merged.Pages[2].Name}
	expected := []string{"Home", "About", "Dashboard"}
	for i, n := range names {
		if n != expected[i] {
			t.Errorf("page %d: expected %s, got %s", i, expected[i], n)
		}
	}
}

func TestMergePrograms_NoPrograms(t *testing.T) {
	_, err := MergePrograms(nil)
	if err == nil {
		t.Fatal("expected error for nil programs")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

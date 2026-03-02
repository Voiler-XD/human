package codegen

import (
	"testing"

	"github.com/barun-bash/human/internal/ir"
)

// testGenerator is a mock CodeGenerator for testing the Registry.
type testGenerator struct {
	name     string
	category Category
	enabled  bool
	stage    string
	outDir   string
}

func (g *testGenerator) Meta() PluginMeta {
	return PluginMeta{
		Name:        g.name,
		Version:     "1.0.0",
		Description: "test generator: " + g.name,
		Category:    g.category,
	}
}

func (g *testGenerator) Enabled(_ *ir.Application) bool { return g.enabled }
func (g *testGenerator) StageName() string               { return g.stage }
func (g *testGenerator) OutputDir() string                { return g.outDir }

func (g *testGenerator) Generate(_ *ir.Application, _ string) error { return nil }

func TestNewRegistryEmpty(t *testing.T) {
	r := NewRegistry()
	if len(r.All()) != 0 {
		t.Errorf("new registry should be empty, got %d generators", len(r.All()))
	}
	if len(r.Names()) != 0 {
		t.Errorf("new registry names should be empty, got %d", len(r.Names()))
	}
}

func TestRegisterAndAll(t *testing.T) {
	r := NewRegistry()
	g1 := &testGenerator{name: "alpha", category: CategoryFrontend, enabled: true, stage: "Generating Alpha"}
	g2 := &testGenerator{name: "beta", category: CategoryBackend, enabled: true, stage: "Generating Beta"}

	if err := r.Register(g1); err != nil {
		t.Fatalf("Register alpha: %v", err)
	}
	if err := r.Register(g2); err != nil {
		t.Fatalf("Register beta: %v", err)
	}

	all := r.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 generators, got %d", len(all))
	}
	if all[0].Meta().Name != "alpha" {
		t.Errorf("first generator = %q, want alpha", all[0].Meta().Name)
	}
	if all[1].Meta().Name != "beta" {
		t.Errorf("second generator = %q, want beta", all[1].Meta().Name)
	}
}

func TestRegisterDuplicate(t *testing.T) {
	r := NewRegistry()
	g1 := &testGenerator{name: "dup", category: CategoryFrontend}
	g2 := &testGenerator{name: "dup", category: CategoryBackend}

	if err := r.Register(g1); err != nil {
		t.Fatalf("first register: %v", err)
	}
	if err := r.Register(g2); err == nil {
		t.Fatal("expected error for duplicate registration")
	}
}

func TestEnabledFiltering(t *testing.T) {
	r := NewRegistry()
	r.Register(&testGenerator{name: "on1", enabled: true, stage: "On1"})
	r.Register(&testGenerator{name: "off", enabled: false, stage: "Off"})
	r.Register(&testGenerator{name: "on2", enabled: true, stage: "On2"})

	app := &ir.Application{Name: "test"}
	enabled := r.Enabled(app)
	if len(enabled) != 2 {
		t.Fatalf("expected 2 enabled generators, got %d", len(enabled))
	}
	if enabled[0].Meta().Name != "on1" {
		t.Errorf("first enabled = %q, want on1", enabled[0].Meta().Name)
	}
	if enabled[1].Meta().Name != "on2" {
		t.Errorf("second enabled = %q, want on2", enabled[1].Meta().Name)
	}
}

func TestGetFound(t *testing.T) {
	r := NewRegistry()
	r.Register(&testGenerator{name: "target", category: CategoryDatabase})

	g := r.Get("target")
	if g == nil {
		t.Fatal("expected to find generator 'target'")
	}
	if g.Meta().Name != "target" {
		t.Errorf("name = %q, want target", g.Meta().Name)
	}
}

func TestGetNotFound(t *testing.T) {
	r := NewRegistry()
	g := r.Get("nonexistent")
	if g != nil {
		t.Errorf("expected nil for nonexistent generator, got %v", g)
	}
}

func TestNamesOrdering(t *testing.T) {
	r := NewRegistry()
	r.Register(&testGenerator{name: "charlie"})
	r.Register(&testGenerator{name: "alpha"})
	r.Register(&testGenerator{name: "bravo"})

	names := r.Names()
	if len(names) != 3 {
		t.Fatalf("expected 3 names, got %d", len(names))
	}
	// Names should be in registration order, not alphabetical.
	if names[0] != "charlie" || names[1] != "alpha" || names[2] != "bravo" {
		t.Errorf("names = %v, want [charlie alpha bravo]", names)
	}
}

func TestPlanStages(t *testing.T) {
	r := NewRegistry()
	r.Register(&testGenerator{name: "a", enabled: true, stage: "Stage A"})
	r.Register(&testGenerator{name: "b", enabled: false, stage: "Stage B"})
	r.Register(&testGenerator{name: "c", enabled: true, stage: "Stage C"})

	app := &ir.Application{Name: "test"}
	stages := r.PlanStages(app)
	if len(stages) != 2 {
		t.Fatalf("expected 2 stages, got %d", len(stages))
	}
	if stages[0] != "Stage A" {
		t.Errorf("stages[0] = %q, want %q", stages[0], "Stage A")
	}
	if stages[1] != "Stage C" {
		t.Errorf("stages[1] = %q, want %q", stages[1], "Stage C")
	}
}

func TestMatchesGoBackend(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"Go", true},
		{"go", true},
		{"Go with Gin", true},
		{"go with fiber", true},
		{"golang", true},
		{"Gin", true},
		{"Node", false},
		{"Python", false},
		{"django", false},
		{"mongodb", false},
		{"", false},
	}

	for _, tt := range tests {
		got := MatchesGoBackend(tt.input)
		if got != tt.want {
			t.Errorf("MatchesGoBackend(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

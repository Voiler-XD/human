package codegen

import (
	"fmt"
	"strings"

	"github.com/barun-bash/human/internal/ir"
)

// Category classifies a code generator by its domain.
type Category string

const (
	CategoryFrontend Category = "frontend"
	CategoryBackend  Category = "backend"
	CategoryDatabase Category = "database"
	CategoryInfra    Category = "infra"
)

// PluginMeta holds descriptive metadata for a code generator.
type PluginMeta struct {
	Name        string   // unique identifier (e.g. "react", "node", "postgres")
	Version     string   // semver (e.g. "1.0.0")
	Description string   // short human-readable description
	Category    Category // frontend, backend, database, or infra
}

// CodeGenerator is the interface that all code generators implement.
// Built-in generators satisfy this interface directly; external plugins
// will be loaded via a future plugin system.
type CodeGenerator interface {
	// Meta returns the generator's metadata.
	Meta() PluginMeta

	// Enabled reports whether this generator should run for the given app.
	// This replaces the hardcoded trigger conditions in pipeline.go.
	Enabled(app *ir.Application) bool

	// StageName returns the display name for progress reporting
	// (e.g. "Generating React frontend").
	StageName() string

	// OutputDir returns the subdirectory name within the build output
	// (e.g. "react", "node"). Empty string means write to the root output dir.
	OutputDir() string

	// Generate writes generated code to outputDir.
	Generate(app *ir.Application, outputDir string) error
}

// Registry holds an ordered collection of code generators.
// Registration order determines execution order.
type Registry struct {
	generators []CodeGenerator
	byName     map[string]CodeGenerator
}

// NewRegistry creates an empty Registry.
func NewRegistry() *Registry {
	return &Registry{
		byName: make(map[string]CodeGenerator),
	}
}

// Register adds a generator to the registry. Returns an error if a generator
// with the same name is already registered.
func (r *Registry) Register(g CodeGenerator) error {
	name := g.Meta().Name
	if _, exists := r.byName[name]; exists {
		return fmt.Errorf("duplicate generator name: %q", name)
	}
	r.generators = append(r.generators, g)
	r.byName[name] = g
	return nil
}

// All returns all registered generators in registration order.
func (r *Registry) All() []CodeGenerator {
	out := make([]CodeGenerator, len(r.generators))
	copy(out, r.generators)
	return out
}

// Enabled returns only the generators that are enabled for the given app,
// preserving registration order.
func (r *Registry) Enabled(app *ir.Application) []CodeGenerator {
	var out []CodeGenerator
	for _, g := range r.generators {
		if g.Enabled(app) {
			out = append(out, g)
		}
	}
	return out
}

// Get returns the generator with the given name, or nil if not found.
func (r *Registry) Get(name string) CodeGenerator {
	return r.byName[name]
}

// Names returns the names of all registered generators in registration order.
func (r *Registry) Names() []string {
	names := make([]string, len(r.generators))
	for i, g := range r.generators {
		names[i] = g.Meta().Name
	}
	return names
}

// PlanStages returns the stage names for all enabled generators, suitable for
// pre-populating a progress display.
func (r *Registry) PlanStages(app *ir.Application) []string {
	var stages []string
	for _, g := range r.generators {
		if g.Enabled(app) {
			stages = append(stages, g.StageName())
		}
	}
	return stages
}

// MatchesGoBackend checks if the backend config indicates Go without
// false-matching strings like "django" or "mongodb".
func MatchesGoBackend(backend string) bool {
	lower := strings.ToLower(backend)
	if lower == "go" || strings.HasPrefix(lower, "go ") {
		return true
	}
	for _, kw := range []string{"gin", "fiber", "golang"} {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

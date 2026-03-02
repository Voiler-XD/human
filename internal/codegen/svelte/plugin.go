package svelte

import (
	"strings"

	"github.com/barun-bash/human/internal/codegen"
	"github.com/barun-bash/human/internal/ir"
)

// Meta returns the generator's metadata.
func (g Generator) Meta() codegen.PluginMeta {
	return codegen.PluginMeta{
		Name:        "svelte",
		Version:     "1.0.0",
		Description: "SvelteKit + TypeScript frontend",
		Category:    codegen.CategoryFrontend,
	}
}

// Enabled reports whether the app's frontend config includes Svelte.
func (g Generator) Enabled(app *ir.Application) bool {
	if app.Config == nil {
		return false
	}
	return strings.Contains(strings.ToLower(app.Config.Frontend), "svelte")
}

// StageName returns the display name for progress reporting.
func (g Generator) StageName() string { return "Generating Svelte frontend" }

// OutputDir returns the subdirectory name within the build output.
func (g Generator) OutputDir() string { return "svelte" }

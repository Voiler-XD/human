package storybook

import (
	"strings"

	"github.com/barun-bash/human/internal/codegen"
	"github.com/barun-bash/human/internal/ir"
)

// Meta returns the generator's metadata.
func (g Generator) Meta() codegen.PluginMeta {
	return codegen.PluginMeta{
		Name:        "storybook",
		Version:     "1.0.0",
		Description: "Storybook stories for frontend components",
		Category:    codegen.CategoryFrontend,
	}
}

// Enabled reports whether a supported frontend framework is configured.
// Only react, vue, angular, and svelte are supported — must match resolveStorybookDir.
func (g Generator) Enabled(app *ir.Application) bool {
	if app.Config == nil {
		return false
	}
	lower := strings.ToLower(app.Config.Frontend)
	return strings.Contains(lower, "react") ||
		strings.Contains(lower, "vue") ||
		strings.Contains(lower, "angular") ||
		strings.Contains(lower, "svelte")
}

// StageName returns the display name for progress reporting.
func (g Generator) StageName() string { return "Generating Storybook stories" }

// OutputDir returns empty because Storybook writes into the frontend directory,
// not a standalone subdirectory. The build pipeline resolves the actual directory.
func (g Generator) OutputDir() string { return "" }

package storybook

import (
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

// Enabled reports whether a frontend framework is configured.
func (g Generator) Enabled(app *ir.Application) bool {
	if app.Config == nil {
		return false
	}
	return app.Config.Frontend != ""
}

// StageName returns the display name for progress reporting.
func (g Generator) StageName() string { return "Generating Storybook stories" }

// OutputDir returns empty because Storybook writes into the frontend directory,
// not a standalone subdirectory. The build pipeline resolves the actual directory.
func (g Generator) OutputDir() string { return "" }

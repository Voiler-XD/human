package architecture

import (
	"github.com/barun-bash/human/internal/codegen"
	"github.com/barun-bash/human/internal/ir"
)

// Meta returns the generator's metadata.
func (g Generator) Meta() codegen.PluginMeta {
	return codegen.PluginMeta{
		Name:        "architecture",
		Version:     "1.0.0",
		Description: "Monolith / microservices / serverless architecture layout",
		Category:    codegen.CategoryInfra,
	}
}

// Enabled reports whether the app has an architecture style configured.
func (g Generator) Enabled(app *ir.Application) bool {
	return app.Architecture != nil && app.Architecture.Style != ""
}

// StageName returns the display name for progress reporting.
func (g Generator) StageName() string { return "Generating architecture layout" }

// OutputDir returns empty because architecture files are written to the root output dir.
func (g Generator) OutputDir() string { return "" }

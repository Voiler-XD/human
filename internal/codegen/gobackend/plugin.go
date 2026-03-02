package gobackend

import (
	"github.com/barun-bash/human/internal/codegen"
	"github.com/barun-bash/human/internal/ir"
)

// Meta returns the generator's metadata.
func (g Generator) Meta() codegen.PluginMeta {
	return codegen.PluginMeta{
		Name:        "go",
		Version:     "1.0.0",
		Description: "Go (Gin/Fiber) backend",
		Category:    codegen.CategoryBackend,
	}
}

// Enabled reports whether the app's backend config indicates Go.
func (g Generator) Enabled(app *ir.Application) bool {
	if app.Config == nil {
		return false
	}
	return codegen.MatchesGoBackend(app.Config.Backend)
}

// StageName returns the display name for progress reporting.
func (g Generator) StageName() string { return "Generating Go backend" }

// OutputDir returns the subdirectory name within the build output.
func (g Generator) OutputDir() string { return "go" }

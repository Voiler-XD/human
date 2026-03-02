package python

import (
	"strings"

	"github.com/barun-bash/human/internal/codegen"
	"github.com/barun-bash/human/internal/ir"
)

// Meta returns the generator's metadata.
func (g Generator) Meta() codegen.PluginMeta {
	return codegen.PluginMeta{
		Name:        "python",
		Version:     "1.0.0",
		Description: "Python (FastAPI/Django) backend",
		Category:    codegen.CategoryBackend,
	}
}

// Enabled reports whether the app's backend config includes Python.
func (g Generator) Enabled(app *ir.Application) bool {
	if app.Config == nil {
		return false
	}
	return strings.Contains(strings.ToLower(app.Config.Backend), "python")
}

// StageName returns the display name for progress reporting.
func (g Generator) StageName() string { return "Generating Python backend" }

// OutputDir returns the subdirectory name within the build output.
func (g Generator) OutputDir() string { return "python" }

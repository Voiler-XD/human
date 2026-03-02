package postgres

import (
	"strings"

	"github.com/barun-bash/human/internal/codegen"
	"github.com/barun-bash/human/internal/ir"
)

// Meta returns the generator's metadata.
func (g Generator) Meta() codegen.PluginMeta {
	return codegen.PluginMeta{
		Name:        "postgres",
		Version:     "1.0.0",
		Description: "PostgreSQL migrations and seeds",
		Category:    codegen.CategoryDatabase,
	}
}

// Enabled reports whether the app's database config includes PostgreSQL.
func (g Generator) Enabled(app *ir.Application) bool {
	if app.Config == nil {
		return false
	}
	return strings.Contains(strings.ToLower(app.Config.Database), "postgres")
}

// StageName returns the display name for progress reporting.
func (g Generator) StageName() string { return "Generating PostgreSQL schema" }

// OutputDir returns the subdirectory name within the build output.
func (g Generator) OutputDir() string { return "postgres" }

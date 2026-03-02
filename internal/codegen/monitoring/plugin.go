package monitoring

import (
	"github.com/barun-bash/human/internal/codegen"
	"github.com/barun-bash/human/internal/ir"
)

// Meta returns the generator's metadata.
func (g Generator) Meta() codegen.PluginMeta {
	return codegen.PluginMeta{
		Name:        "monitoring",
		Version:     "1.0.0",
		Description: "Prometheus + Grafana monitoring configuration",
		Category:    codegen.CategoryInfra,
	}
}

// Enabled reports whether the app has monitoring rules configured.
func (g Generator) Enabled(app *ir.Application) bool {
	return len(app.Monitoring) > 0
}

// StageName returns the display name for progress reporting.
func (g Generator) StageName() string { return "Generating monitoring configuration" }

// OutputDir returns the subdirectory name within the build output.
func (g Generator) OutputDir() string { return "monitoring" }

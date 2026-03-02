package docker

import (
	"strings"

	"github.com/barun-bash/human/internal/codegen"
	"github.com/barun-bash/human/internal/ir"
)

// Meta returns the generator's metadata.
func (g Generator) Meta() codegen.PluginMeta {
	return codegen.PluginMeta{
		Name:        "docker",
		Version:     "1.0.0",
		Description: "Dockerfile and docker-compose configuration",
		Category:    codegen.CategoryInfra,
	}
}

// Enabled reports whether the app's deploy config includes Docker.
func (g Generator) Enabled(app *ir.Application) bool {
	if app.Config == nil {
		return false
	}
	return strings.Contains(strings.ToLower(app.Config.Deploy), "docker")
}

// StageName returns the display name for progress reporting.
func (g Generator) StageName() string { return "Generating Docker configuration" }

// OutputDir returns empty because Docker files are written to the root output dir.
func (g Generator) OutputDir() string { return "" }

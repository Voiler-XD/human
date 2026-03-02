package cicd

import (
	"github.com/barun-bash/human/internal/codegen"
	"github.com/barun-bash/human/internal/ir"
)

// Meta returns the generator's metadata.
func (g Generator) Meta() codegen.PluginMeta {
	return codegen.PluginMeta{
		Name:        "cicd",
		Version:     "1.0.0",
		Description: "GitHub Actions CI/CD workflows",
		Category:    codegen.CategoryInfra,
	}
}

// Enabled always returns true — CI/CD pipelines are generated for every app.
func (g Generator) Enabled(_ *ir.Application) bool { return true }

// StageName returns the display name for progress reporting.
func (g Generator) StageName() string { return "Generating CI/CD pipelines" }

// OutputDir returns empty because CI/CD files are written to the root output dir.
func (g Generator) OutputDir() string { return "" }

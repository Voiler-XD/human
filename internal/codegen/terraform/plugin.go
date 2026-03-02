package terraform

import (
	"strings"

	"github.com/barun-bash/human/internal/codegen"
	"github.com/barun-bash/human/internal/ir"
)

// Meta returns the generator's metadata.
func (g Generator) Meta() codegen.PluginMeta {
	return codegen.PluginMeta{
		Name:        "terraform",
		Version:     "1.0.0",
		Description: "Terraform infrastructure (AWS ECS/RDS, GCP Cloud Run/SQL)",
		Category:    codegen.CategoryInfra,
	}
}

// Enabled reports whether the app's deploy config includes a cloud provider or Terraform.
func (g Generator) Enabled(app *ir.Application) bool {
	if app.Config == nil {
		return false
	}
	d := strings.ToLower(app.Config.Deploy)
	return strings.Contains(d, "aws") || strings.Contains(d, "gcp") || strings.Contains(d, "terraform")
}

// StageName returns the display name for progress reporting.
func (g Generator) StageName() string { return "Generating Terraform infrastructure" }

// OutputDir returns the subdirectory name within the build output.
func (g Generator) OutputDir() string { return "terraform" }

package build

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/barun-bash/human/internal/codegen"
	"github.com/barun-bash/human/internal/codegen/angular"
	"github.com/barun-bash/human/internal/codegen/architecture"
	"github.com/barun-bash/human/internal/codegen/cicd"
	"github.com/barun-bash/human/internal/codegen/docker"
	"github.com/barun-bash/human/internal/codegen/gobackend"
	"github.com/barun-bash/human/internal/codegen/monitoring"
	"github.com/barun-bash/human/internal/codegen/node"
	"github.com/barun-bash/human/internal/codegen/postgres"
	"github.com/barun-bash/human/internal/codegen/python"
	"github.com/barun-bash/human/internal/codegen/react"
	"github.com/barun-bash/human/internal/codegen/storybook"
	"github.com/barun-bash/human/internal/codegen/svelte"
	"github.com/barun-bash/human/internal/codegen/terraform"
	"github.com/barun-bash/human/internal/codegen/vue"
	"github.com/barun-bash/human/internal/ir"
)

// DefaultRegistry returns a registry populated with all 14 built-in code
// generators in the correct execution order. Quality and scaffold are NOT
// included — they are run as explicit post-loop steps in the pipeline.
func DefaultRegistry() *codegen.Registry {
	reg := codegen.NewRegistry()

	// Registration order = execution order.
	// Frontend generators first, then backend, database, infrastructure.
	generators := []codegen.CodeGenerator{
		react.Generator{},
		vue.Generator{},
		angular.Generator{},
		svelte.Generator{},
		storybook.Generator{},
		node.Generator{},
		python.Generator{},
		gobackend.Generator{},
		postgres.Generator{},
		docker.Generator{},
		cicd.Generator{},
		terraform.Generator{},
		architecture.Generator{},
		monitoring.Generator{},
	}

	for _, g := range generators {
		// Built-in generators should never have duplicate names.
		if err := reg.Register(g); err != nil {
			panic("built-in generator registration: " + err.Error())
		}
	}

	return reg
}

// resolveStorybookDir determines the frontend output directory for Storybook.
// Storybook generates into the frontend directory, not a standalone subdirectory.
func resolveStorybookDir(app *ir.Application, outputDir string) string {
	if app.Config == nil {
		return ""
	}
	frontendLower := strings.ToLower(app.Config.Frontend)
	switch {
	case strings.Contains(frontendLower, "react"):
		return filepath.Join(outputDir, "react")
	case strings.Contains(frontendLower, "vue"):
		return filepath.Join(outputDir, "vue")
	case strings.Contains(frontendLower, "angular"):
		return filepath.Join(outputDir, "angular")
	case strings.Contains(frontendLower, "svelte"):
		return filepath.Join(outputDir, "svelte")
	}
	return ""
}

// countStorybookFiles counts the Storybook-specific files generated in a
// frontend directory.
func countStorybookFiles(frontendDir string) int {
	return countFilesUnder(filepath.Join(frontendDir, ".storybook")) +
		countFilesUnder(filepath.Join(frontendDir, "src", "stories"))
}

// countScaffoldFiles counts the scaffold-generated files across the output.
func countScaffoldFiles(outputDir string) int {
	count := 0
	for _, name := range []string{"package.json", "README.md", ".env.example", "start.sh"} {
		if _, err := os.Stat(filepath.Join(outputDir, name)); err == nil {
			count++
		}
	}
	for _, sub := range []string{"node", "react", "vue"} {
		for _, name := range []string{"package.json", "tsconfig.json", "vite.config.ts"} {
			if _, err := os.Stat(filepath.Join(outputDir, sub, name)); err == nil {
				count++
			}
		}
	}
	return count
}

// countFilesUnder returns the number of regular files under dir.
func countFilesUnder(dir string) int {
	count := 0
	filepath.WalkDir(dir, func(_ string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			count++
		}
		return nil
	})
	return count
}

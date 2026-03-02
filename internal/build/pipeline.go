package build

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/barun-bash/human/internal/codegen"
	"github.com/barun-bash/human/internal/codegen/scaffold"
	"github.com/barun-bash/human/internal/codegen/storybook"
	"github.com/barun-bash/human/internal/ir"
	"github.com/barun-bash/human/internal/quality"
)

// Result tracks the output of a single generator.
type Result struct {
	Name     string
	Dir      string
	Files    int
	Duration time.Duration
}

// BuildTiming holds the total build duration.
type BuildTiming struct {
	Total time.Duration
}

// MatchesGoBackend checks if the backend config indicates Go without
// false-matching strings like "django" or "mongodb".
func MatchesGoBackend(backend string) bool {
	return codegen.MatchesGoBackend(backend)
}

// CountFiles returns the number of regular files under dir.
func CountFiles(dir string) int {
	count := 0
	filepath.WalkDir(dir, func(_ string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			count++
		}
		return nil
	})
	return count
}

// PlanStages returns the list of stage names that will run for the given app.
// Use this to pre-populate a progress display.
func PlanStages(app *ir.Application) []string {
	return PlanStagesWithRegistry(DefaultRegistry(), app)
}

// PlanStagesWithRegistry returns the list of stage names for the given registry
// and app. Includes quality and scaffold stages that always run.
func PlanStagesWithRegistry(reg *codegen.Registry, app *ir.Application) []string {
	stages := reg.PlanStages(app)
	stages = append(stages, "Running quality checks")
	stages = append(stages, "Scaffolding project files")
	return stages
}

// ProgressFunc is called before each build stage with the stage name.
type ProgressFunc func(stage string)

// RunGenerators dispatches all code generators based on the app's build config,
// then runs the quality engine and scaffolder. Returns build results for each
// generator, the quality result, and build timing.
func RunGenerators(app *ir.Application, outputDir string) ([]Result, *quality.Result, *BuildTiming, error) {
	return RunGeneratorsWithProgress(app, outputDir, nil)
}

// RunGeneratorsWithProgress is like RunGenerators but calls progress before each stage.
func RunGeneratorsWithProgress(app *ir.Application, outputDir string, progress ProgressFunc) ([]Result, *quality.Result, *BuildTiming, error) {
	return RunGeneratorsWithRegistry(DefaultRegistry(), app, outputDir, progress)
}

// RunGeneratorsWithRegistry dispatches generators from the given registry,
// then runs the quality engine and scaffolder. This allows custom registries
// for testing or plugin scenarios.
func RunGeneratorsWithRegistry(reg *codegen.Registry, app *ir.Application, outputDir string, progress ProgressFunc) ([]Result, *quality.Result, *BuildTiming, error) {
	buildStart := time.Now()
	var results []Result

	report := func(stage string) {
		if progress != nil {
			progress(stage)
		}
	}

	timeGen := func(name, dir string, files int, start time.Time) Result {
		return Result{Name: name, Dir: dir, Files: files, Duration: time.Since(start)}
	}

	// Run all enabled generators from the registry.
	for _, g := range reg.Enabled(app) {
		name := g.Meta().Name
		report(g.StageName())
		start := time.Now()

		// Resolve target directory.
		var dir string
		switch name {
		case "storybook":
			// Storybook generates into the frontend directory, not standalone.
			dir = resolveStorybookDir(app, outputDir)
			if dir == "" {
				continue
			}
		default:
			if od := g.OutputDir(); od != "" {
				dir = filepath.Join(outputDir, od)
			} else {
				dir = outputDir
			}
		}

		// For Docker, count files before generation so we can diff.
		var beforeCount int
		if name == "docker" {
			beforeCount = CountFiles(outputDir)
		}

		// Run the generator.
		if err := g.Generate(app, dir); err != nil {
			return nil, nil, nil, fmt.Errorf("%s codegen: %w", name, err)
		}

		// Storybook: call GetFramework so scaffold knows what deps to inject.
		if name == "storybook" {
			_ = storybook.GetFramework(app)
		}

		// Count generated files — each generator has a different counting strategy.
		var files int
		switch name {
		case "storybook":
			files = countStorybookFiles(dir)
		case "docker":
			files = CountFiles(outputDir) - beforeCount
		case "cicd":
			files = CountFiles(filepath.Join(outputDir, ".github"))
		case "architecture":
			files = CountFiles(filepath.Join(outputDir, "services")) +
				CountFiles(filepath.Join(outputDir, "functions")) +
				CountFiles(filepath.Join(outputDir, "gateway"))
			if files == 0 {
				continue // architecture generator produced nothing (monolith)
			}
		default:
			files = CountFiles(dir)
		}

		results = append(results, timeGen(name, dir, files, start))
	}

	// Quality engine — always runs after code generators.
	report("Running quality checks")
	qualityStart := time.Now()
	qResult, err := quality.Run(app, outputDir)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("quality engine: %w", err)
	}
	qualityFiles := qResult.TestFiles + qResult.ComponentTestFiles + qResult.EdgeTestFiles + 3
	results = append(results, timeGen("quality", outputDir, qualityFiles, qualityStart))

	// Scaffolder — always runs last.
	report("Scaffolding project files")
	scaffoldStart := time.Now()
	sg := scaffold.Generator{}
	if err := sg.Generate(app, outputDir); err != nil {
		return nil, nil, nil, fmt.Errorf("scaffold: %w", err)
	}
	results = append(results, timeGen("scaffold", outputDir, countScaffoldFiles(outputDir), scaffoldStart))

	timing := &BuildTiming{Total: time.Since(buildStart)}
	return results, qResult, timing, nil
}

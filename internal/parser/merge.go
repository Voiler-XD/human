package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// DiscoverFiles finds all .human files for a project given a path.
// If path is a .human file, it discovers sibling .human files in the same directory.
// If path is a directory, it finds all .human files in it.
// Returns files ordered: app.human first (if present), then alphabetical.
// When multiple files exist, app.human is required — returns an error if missing.
func DiscoverFiles(path string) ([]string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("cannot access %s: %w", path, err)
	}

	var dir string
	if info.IsDir() {
		dir = path
	} else {
		dir = filepath.Dir(path)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading directory %s: %w", dir, err)
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.HasSuffix(e.Name(), ".human") {
			files = append(files, filepath.Join(dir, e.Name()))
		}
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no .human files found in %s", dir)
	}

	// Single file — return as-is (backward compatible, no app.human requirement).
	if len(files) == 1 {
		return files, nil
	}

	// Multiple files — app.human is required.
	hasApp := false
	for _, f := range files {
		if filepath.Base(f) == "app.human" {
			hasApp = true
			break
		}
	}
	if !hasApp {
		return nil, fmt.Errorf("multi-file project requires app.human in %s", dir)
	}

	// Sort: app.human first, then alphabetical.
	sort.Slice(files, func(i, j int) bool {
		bi, bj := filepath.Base(files[i]), filepath.Base(files[j])
		if bi == "app.human" {
			return true
		}
		if bj == "app.human" {
			return false
		}
		return bi < bj
	})

	return files, nil
}

// ParseFiles parses each file and returns the resulting programs.
// The File field is set on every declaration in each program.
// Fails fast on the first parse error, including the filename in the message.
func ParseFiles(files []string) ([]*Program, error) {
	programs := make([]*Program, 0, len(files))

	for _, file := range files {
		source, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", file, err)
		}

		prog, err := Parse(string(source))
		if err != nil {
			return nil, fmt.Errorf("parse error in %s: %w", file, err)
		}

		// Tag every declaration with its source file.
		tagFile(prog, file)
		programs = append(programs, prog)
	}

	return programs, nil
}

// tagFile sets the File field on all declarations in a program.
func tagFile(prog *Program, file string) {
	if prog.App != nil {
		prog.App.File = file
	}
	for _, d := range prog.Data {
		d.File = file
	}
	for _, d := range prog.Pages {
		d.File = file
	}
	for _, d := range prog.Components {
		d.File = file
	}
	for _, d := range prog.APIs {
		d.File = file
	}
	for _, d := range prog.Policies {
		d.File = file
	}
	for _, d := range prog.Workflows {
		d.File = file
	}
	if prog.Theme != nil {
		prog.Theme.File = file
	}
	if prog.Authentication != nil {
		prog.Authentication.File = file
	}
	if prog.Database != nil {
		prog.Database.File = file
	}
	for _, d := range prog.Integrations {
		d.File = file
	}
	for _, d := range prog.Environments {
		d.File = file
	}
	for _, d := range prog.ErrorHandlers {
		d.File = file
	}
	if prog.Build != nil {
		prog.Build.File = file
	}
	if prog.Architecture != nil {
		prog.Architecture.File = file
	}
}

// MergePrograms combines multiple parsed programs into a single program.
// Singleton declarations (App, Theme, Authentication, Database, Build, Architecture)
// use the first non-nil value; duplicates produce an error with both filenames.
// Slice declarations are appended in file order.
func MergePrograms(programs []*Program) (*Program, error) {
	if len(programs) == 0 {
		return nil, fmt.Errorf("no programs to merge")
	}
	if len(programs) == 1 {
		return programs[0], nil
	}

	merged := &Program{}

	for _, prog := range programs {
		// Singleton: App
		if prog.App != nil {
			if merged.App != nil {
				return nil, fmt.Errorf("duplicate app declaration: %s (line %d) and %s (line %d)",
					merged.App.File, merged.App.Line, prog.App.File, prog.App.Line)
			}
			merged.App = prog.App
		}

		// Singleton: Theme
		if prog.Theme != nil {
			if merged.Theme != nil {
				return nil, fmt.Errorf("duplicate theme declaration: %s (line %d) and %s (line %d)",
					merged.Theme.File, merged.Theme.Line, prog.Theme.File, prog.Theme.Line)
			}
			merged.Theme = prog.Theme
		}

		// Singleton: Authentication
		if prog.Authentication != nil {
			if merged.Authentication != nil {
				return nil, fmt.Errorf("duplicate authentication declaration: %s (line %d) and %s (line %d)",
					merged.Authentication.File, merged.Authentication.Line, prog.Authentication.File, prog.Authentication.Line)
			}
			merged.Authentication = prog.Authentication
		}

		// Singleton: Database
		if prog.Database != nil {
			if merged.Database != nil {
				return nil, fmt.Errorf("duplicate database declaration: %s (line %d) and %s (line %d)",
					merged.Database.File, merged.Database.Line, prog.Database.File, prog.Database.Line)
			}
			merged.Database = prog.Database
		}

		// Singleton: Build
		if prog.Build != nil {
			if merged.Build != nil {
				return nil, fmt.Errorf("duplicate build declaration: %s (line %d) and %s (line %d)",
					merged.Build.File, merged.Build.Line, prog.Build.File, prog.Build.Line)
			}
			merged.Build = prog.Build
		}

		// Singleton: Architecture
		if prog.Architecture != nil {
			if merged.Architecture != nil {
				return nil, fmt.Errorf("duplicate architecture declaration: %s (line %d) and %s (line %d)",
					merged.Architecture.File, merged.Architecture.Line, prog.Architecture.File, prog.Architecture.Line)
			}
			merged.Architecture = prog.Architecture
		}

		// Slices: append in file order
		merged.Data = append(merged.Data, prog.Data...)
		merged.Pages = append(merged.Pages, prog.Pages...)
		merged.Components = append(merged.Components, prog.Components...)
		merged.APIs = append(merged.APIs, prog.APIs...)
		merged.Policies = append(merged.Policies, prog.Policies...)
		merged.Workflows = append(merged.Workflows, prog.Workflows...)
		merged.Integrations = append(merged.Integrations, prog.Integrations...)
		merged.Environments = append(merged.Environments, prog.Environments...)
		merged.ErrorHandlers = append(merged.ErrorHandlers, prog.ErrorHandlers...)

		// Sections and Statements: append in order
		merged.Sections = append(merged.Sections, prog.Sections...)
		merged.Statements = append(merged.Statements, prog.Statements...)
	}

	return merged, nil
}

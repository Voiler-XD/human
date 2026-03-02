# Building a Code Generator for Human

This guide explains how to add a new code generator to the Human compiler. A new generator could be a frontend framework (e.g., HTMX), a backend (e.g., Rust+Axum), a database (e.g., MongoDB), or infrastructure target (e.g., Kubernetes).

## Architecture Overview

The Human compiler works in three stages:

```
.human file → Lexer → Parser → Intent IR → Code Generators → Output
```

Code generators receive an `*ir.Application` (the Intent IR) and produce files. They never touch the parser or lexer. This means you can add a new target without understanding the Human language grammar — you only need to understand the IR types.

## Quick Start

### 1. Create your generator package

```
internal/codegen/
├── react/         ← existing frontend
├── node/          ← existing backend
├── postgres/      ← existing database
├── yourtarget/    ← create this
│   ├── generator.go
│   ├── generator_test.go
│   └── helpers.go       (optional)
```

### 2. Implement the CodeGenerator interface

All generators implement the `codegen.CodeGenerator` interface defined in `internal/codegen/plugin.go`. This interface has 5 methods:

{% raw %}
```go
// CodeGenerator is the interface that all code generators implement.
type CodeGenerator interface {
    Meta() PluginMeta                                    // metadata (name, version, category)
    Enabled(app *ir.Application) bool                    // should this generator run?
    StageName() string                                   // progress display name
    OutputDir() string                                   // subdirectory (empty = root)
    Generate(app *ir.Application, outputDir string) error // produce files
}

// PluginMeta holds descriptive metadata for a code generator.
type PluginMeta struct {
    Name        string   // unique identifier (e.g. "react", "node")
    Version     string   // semver (e.g. "1.0.0")
    Description string   // short human-readable description
    Category    Category // "frontend", "backend", "database", or "infra"
}
```
{% endraw %}

Create two files in your package — `generator.go` for the `Generate` method and `plugin.go` for the interface methods:

{% raw %}
```go
// internal/codegen/yourtarget/generator.go
package yourtarget

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/barun-bash/human/internal/ir"
)

type Generator struct{}

func (g Generator) Generate(app *ir.Application, outputDir string) error {
    dirs := []string{
        filepath.Join(outputDir, "src"),
        filepath.Join(outputDir, "src", "models"),
    }
    for _, d := range dirs {
        if err := os.MkdirAll(d, 0755); err != nil {
            return fmt.Errorf("creating directory %s: %w", d, err)
        }
    }

    files := map[string]string{
        filepath.Join(outputDir, "src", "main.go"): generateMain(app),
    }
    for path, content := range files {
        if err := writeFile(path, content); err != nil {
            return err
        }
    }
    return nil
}
```
{% endraw %}

{% raw %}
```go
// internal/codegen/yourtarget/plugin.go
package yourtarget

import (
    "strings"

    "github.com/barun-bash/human/internal/codegen"
    "github.com/barun-bash/human/internal/ir"
)

func (g Generator) Meta() codegen.PluginMeta {
    return codegen.PluginMeta{
        Name:        "yourtarget",
        Version:     "1.0.0",
        Description: "YourTarget backend",
        Category:    codegen.CategoryBackend,
    }
}

func (g Generator) Enabled(app *ir.Application) bool {
    if app.Config == nil {
        return false
    }
    return strings.Contains(strings.ToLower(app.Config.Backend), "yourtarget")
}

func (g Generator) StageName() string { return "Generating YourTarget backend" }
func (g Generator) OutputDir() string { return "yourtarget" }
```
{% endraw %}

### 3. Register in the build pipeline

Generators are registered in `internal/build/registry.go` via the `DefaultRegistry()` function. The registry determines execution order and dispatches only enabled generators.

Add your import and registration:

```go
// In internal/build/registry.go
import "github.com/barun-bash/human/internal/codegen/yourtarget"

func DefaultRegistry() *codegen.Registry {
    reg := codegen.NewRegistry()
    // ... existing generators ...
    reg.Register(yourtarget.Generator{})
    // ...
    return reg
}
```

Registration order determines execution order. Place your generator in the appropriate section (frontend, backend, database, or infrastructure). The `Enabled` method on your generator controls when it runs — no `if` blocks needed in the pipeline.

## IR Type Reference

Every generator reads from `*ir.Application`. Here are the types you need to handle:

### Core Types

| IR Type | Field | What to generate |
|---------|-------|-----------------|
| `app.Name` | `string` | Project name, package name |
| `app.Platform` | `string` | "web", "mobile", "api" |
| `app.Config` | `*BuildConfig` | Frontend/backend/database/deploy targets |
| `app.Data` | `[]*DataModel` | Database models, ORM classes, type definitions |
| `app.Pages` | `[]*Page` | UI components, templates, routes |
| `app.Components` | `[]*Component` | Reusable UI components |
| `app.APIs` | `[]*Endpoint` | Route handlers, controllers |
| `app.Auth` | `*Auth` | Authentication middleware, login/signup |
| `app.Policies` | `[]*Policy` | Authorization rules, RBAC |
| `app.Theme` | `*Theme` | CSS variables, design tokens, colors |
| `app.Workflows` | `[]*Workflow` | Event handlers, background jobs |
| `app.Integrations` | `[]*Integration` | Third-party services (Stripe, SendGrid, S3) |
| `app.Database` | `*DatabaseConfig` | Engine, indexes, seeding rules |
| `app.Architecture` | `*Architecture` | Monolith / microservices / serverless |
| `app.Monitoring` | `[]*MonitoringRule` | Prometheus alerting rules |
| `app.Environments` | `[]*Environment` | Staging/production config |
| `app.ErrorHandlers` | `[]*ErrorHandler` | Custom error pages, handlers |

### DataModel

```go
type DataModel struct {
    Name      string
    Fields    []*DataField
    Relations []*Relation
}

type DataField struct {
    Name       string
    Type       string   // "text", "number", "email", "date", "datetime", "boolean", "json", etc.
    Required   bool
    Unique     bool
    Default    string
    Enum       []string // non-nil if the field is an enum
    Validation []*ValidationRule
}

type Relation struct {
    Type   string // "has_many", "has_one", "belongs_to"
    Target string // name of related model
    Field  string
}
```

### Page and Action

Pages contain `[]*Action` — the most complex IR type. Actions represent English statements about what a page does:

```go
type Action struct {
    Type       string // see table below
    Text       string // original English text
    Target     string // model or page name if detected
    Field      string
    Condition  string
    Children   []*Action // nested actions (e.g., loop body, conditional branches)
}
```

| Action Type | Meaning | Example generation |
|-------------|---------|-------------------|
| `display` | Show data | Render list/detail component |
| `interact` | User interaction | onClick/onSubmit handler |
| `input` | Form input | Input field with validation |
| `navigate` | Page navigation | Router link / redirect |
| `query` | Data fetching | API call, database query |
| `create` | Create record | POST request, form submission |
| `update` | Update record | PUT/PATCH request |
| `delete` | Delete record | DELETE request with confirmation |
| `loop` | Iteration | `.map()` / `v-for` / `*ngFor` |
| `condition` | Conditional | `if`/`else` render, empty state |
| `validate` | Validation | Client-side validation check |
| `compute` | Calculation | Derived value, computed property |
| `filter` | Filtering | Search bar, filter dropdown |
| `sort` | Sorting | Sort controls |
| `paginate` | Pagination | Page controls, infinite scroll |
| `upload` | File upload | File input, drag & drop |
| `notify` | Notification | Toast, alert, confirmation |
| `subscribe` | Real-time | WebSocket, SSE listener |
| `authenticate` | Auth action | Login/logout/signup handler |

Study `internal/codegen/react/pages.go` for the most complete example of Action interpretation.

### Endpoint

```go
type Endpoint struct {
    Name       string
    Method     string // "GET", "POST", "PUT", "DELETE"
    Path       string
    Auth       bool
    Params     []*Param
    Steps      []*Action
    Validation []*ValidationRule
}
```

### Theme

```go
type Theme struct {
    DesignSystem string // "Material", "Shadcn", "Tailwind", etc.
    Colors       map[string]string
    Fonts        map[string]string
    Spacing      map[string]string
    BorderRadius map[string]string
}
```

Theme integration is framework-specific. The `internal/codegen/themes/` package provides helpers:
- `themes.GenerateReactTheme(theme)` → CSS variables + React theme files
- `themes.GenerateVueTheme(theme)` → Vue theme files
- If you add a new frontend, add a corresponding theme helper.

For the complete IR type reference, see [IR_SCHEMA.md](IR_SCHEMA.md).

## Existing Generators to Study

| Generator | Best example of |
|-----------|----------------|
| `react/` | Page generation, Action interpretation, theme integration, component generation |
| `vue/` | Alternative frontend showing same IR → different framework |
| `node/` | API routes, middleware, validation, policies, Prisma schema |
| `gobackend/` | Go backend (Gin), shows same API pattern in a different language |
| `python/` | FastAPI/Django backend, different ORM patterns |
| `postgres/` | SQL migrations, indexes, seed data |
| `docker/` | Infrastructure config, multi-stage builds |
| `terraform/` | Cloud infrastructure (AWS ECS/RDS, GCP Cloud Run) |
| `themes/` | Design system integration, CSS variables |
| `storybook/` | Generates into an existing frontend directory (not standalone) |
| `scaffold/` | package.json, README, start scripts — always runs last |

## File Writing Pattern

All generators use the same pattern: build a `map[string]string` of path → content, then write:

```go
files := map[string]string{
    filepath.Join(outputDir, "src", "main.ts"): content1,
    filepath.Join(outputDir, "config.json"):    content2,
}
for path, content := range files {
    if err := writeFile(path, content); err != nil {
        return err
    }
}
```

For per-model or per-page files, append to the map in a loop:

```go
for _, model := range app.Data {
    filename := toSnakeCase(model.Name) + ".ts"
    path := filepath.Join(outputDir, "src", "models", filename)
    files[path] = generateModel(model)
}
```

## Testing

Every generator should have a `_test.go` file. Use stdlib `testing` only (no external test libraries).

{% raw %}
```go
package yourtarget

import (
    "os"
    "path/filepath"
    "strings"
    "testing"

    "github.com/barun-bash/human/internal/ir"
)

func TestGenerate(t *testing.T) {
    app := &ir.Application{
        Name:     "TestApp",
        Platform: "web",
        Config: &ir.BuildConfig{
            Frontend: "YourTarget",
            Backend:  "Node",
            Database: "PostgreSQL",
        },
        Data: []*ir.DataModel{{
            Name: "Task",
            Fields: []*ir.DataField{
                {Name: "title", Type: "text", Required: true},
                {Name: "status", Type: "text", Enum: []string{"pending", "done"}},
            },
        }},
        Pages: []*ir.Page{{
            Name: "Tasks",
            Content: []*ir.Action{
                {Type: "display", Text: "show a list of tasks", Target: "Task"},
            },
        }},
        APIs: []*ir.Endpoint{{
            Name:   "GetTasks",
            Method: "GET",
            Path:   "/tasks",
        }},
    }

    dir := t.TempDir()
    err := Generator{}.Generate(app, dir)
    if err != nil {
        t.Fatalf("Generate failed: %v", err)
    }

    // Verify expected files exist
    expected := []string{
        "src/main.go",
        "src/models/task.go",
    }
    for _, rel := range expected {
        path := filepath.Join(dir, rel)
        if _, err := os.Stat(path); os.IsNotExist(err) {
            t.Errorf("expected file %s not found", rel)
        }
    }

    // Verify file contents contain key patterns
    content, _ := os.ReadFile(filepath.Join(dir, "src/models/task.go"))
    if !strings.Contains(string(content), "Task") {
        t.Error("model file should contain Task")
    }
}
```
{% endraw %}

Run tests with:
```bash
go test ./internal/codegen/yourtarget/...
```

## Plugin Configuration

Generators can be overridden per-project via `.human/config.json`:

```json
{
  "plugins": [
    {"name": "react", "enabled": true},
    {"name": "storybook", "enabled": false},
    {"name": "docker", "settings": {"compose_version": "3.8"}}
  ]
}
```

- **`enabled: true`** — force the generator to run even if `Enabled()` returns false
- **`enabled: false`** — force the generator to skip even if `Enabled()` returns true
- **`enabled` omitted** — use the generator's own `Enabled()` logic (default)
- **`settings`** — arbitrary key-value pairs passed to the generator

## Checklist

Before submitting a new generator:

- [ ] Package created at `internal/codegen/yourtarget/`
- [ ] `Generator` struct implements `codegen.CodeGenerator` interface (5 methods)
- [ ] `plugin.go` with `Meta()`, `Enabled()`, `StageName()`, `OutputDir()`
- [ ] `generator.go` with `Generate(app *ir.Application, outputDir string) error`
- [ ] Handles `app.Data` → models/types
- [ ] Handles `app.Pages` → UI components (if frontend)
- [ ] Handles `app.APIs` → routes/controllers (if backend)
- [ ] Handles `app.Auth` → authentication middleware
- [ ] Handles `app.Theme` → design tokens/CSS (if frontend)
- [ ] Tests with a representative `*ir.Application`
- [ ] Registered in `internal/build/registry.go` (`DefaultRegistry()`)
- [ ] `go vet ./...` passes
- [ ] `go test ./...` passes
- [ ] Generates working, runnable code (not stubs)

## Dependencies

The Human compiler uses **zero external Go dependencies** (only `golang.org/x/sys` and `golang.org/x/term` for terminal support). All code generation uses `fmt.Sprintf` and `strings.Builder` — no template engines. Keep it that way.

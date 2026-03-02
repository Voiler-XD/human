package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Simplified IR types — only the fields this generator needs.
type application struct {
	Name   string      `json:"name"`
	Config *buildCfg   `json:"config,omitempty"`
	Pages  []*page     `json:"pages,omitempty"`
	Theme  *theme      `json:"theme,omitempty"`
	Data   []*model    `json:"data,omitempty"`
	APIs   []*endpoint `json:"apis,omitempty"`
}

type buildCfg struct {
	Frontend string `json:"frontend"`
	Backend  string `json:"backend"`
}

type page struct {
	Name    string    `json:"name"`
	Route   string    `json:"route"`
	Content []*action `json:"content,omitempty"`
}

type action struct {
	Type     string    `json:"type"`
	Text     string    `json:"text"`
	Target   string    `json:"target"`
	Children []*action `json:"children,omitempty"`
}

type theme struct {
	DesignSystem string            `json:"design_system"`
	Colors       map[string]string `json:"colors,omitempty"`
	Fonts        map[string]string `json:"fonts,omitempty"`
}

type model struct {
	Name   string  `json:"name"`
	Fields []*field `json:"fields,omitempty"`
}

type field struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
}

type endpoint struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Name   string `json:"name"`
}

func runGenerate(args []string) error {
	fs := flag.NewFlagSet("generate", flag.ExitOnError)
	irFile := fs.String("ir", "", "Path to IR JSON file")
	outputDir := fs.String("output", "", "Output directory")
	_ = fs.String("settings", "", "Plugin settings JSON")
	fs.Parse(args)

	if *irFile == "" || *outputDir == "" {
		return fmt.Errorf("--ir and --output are required")
	}

	data, err := os.ReadFile(*irFile)
	if err != nil {
		return fmt.Errorf("reading IR: %w", err)
	}

	var app application
	if err := json.Unmarshal(data, &app); err != nil {
		return fmt.Errorf("parsing IR: %w", err)
	}

	return generate(&app, *outputDir)
}

func generate(app *application, outputDir string) error {
	dirs := []string{
		outputDir,
		filepath.Join(outputDir, "templates"),
		filepath.Join(outputDir, "static"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return err
		}
	}

	files := map[string]string{
		filepath.Join(outputDir, "templates", "base.html"):  generateBaseTemplate(app),
		filepath.Join(outputDir, "static", "app.js"):        generateAlpineJS(app),
		filepath.Join(outputDir, "static", "styles.css"):    generateStyles(app),
	}

	// Generate a template for each page.
	for _, p := range app.Pages {
		filename := strings.ToLower(p.Name) + ".html"
		path := filepath.Join(outputDir, "templates", filename)
		files[path] = generatePageTemplate(app, p)
	}

	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

func generateBaseTemplate(app *application) string {
	var b strings.Builder
	b.WriteString("<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n")
	b.WriteString("  <meta charset=\"UTF-8\">\n")
	b.WriteString("  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n")
	b.WriteString(fmt.Sprintf("  <title>%s</title>\n", app.Name))
	b.WriteString("  <script src=\"https://unpkg.com/htmx.org@2.0.0\"></script>\n")
	b.WriteString("  <script defer src=\"https://unpkg.com/alpinejs@3.x.x/dist/cdn.min.js\"></script>\n")
	b.WriteString("  <link rel=\"stylesheet\" href=\"/static/styles.css\">\n")
	b.WriteString("</head>\n<body>\n")
	b.WriteString("  <nav>\n")
	for _, p := range app.Pages {
		route := p.Route
		if route == "" {
			route = "/" + strings.ToLower(p.Name)
		}
		b.WriteString(fmt.Sprintf("    <a href=\"%s\" hx-get=\"%s\" hx-target=\"#content\" hx-push-url=\"true\">%s</a>\n", route, route, p.Name))
	}
	b.WriteString("  </nav>\n")
	b.WriteString("  <main id=\"content\">\n")
	b.WriteString("    {{block \"content\" .}}{{end}}\n")
	b.WriteString("  </main>\n")
	b.WriteString("  <script src=\"/static/app.js\"></script>\n")
	b.WriteString("</body>\n</html>\n")
	return b.String()
}

func generatePageTemplate(app *application, p *page) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("{{define \"content\"}}\n<div class=\"page\" id=\"page-%s\">\n", strings.ToLower(p.Name)))
	b.WriteString(fmt.Sprintf("  <h1>%s</h1>\n", p.Name))

	for _, a := range p.Content {
		switch a.Type {
		case "display":
			target := a.Target
			if target == "" {
				target = "items"
			}
			apiPath := "/api/" + strings.ToLower(target) + "s"
			b.WriteString(fmt.Sprintf("  <div hx-get=\"%s\" hx-trigger=\"load\" hx-target=\"#%s-list\">\n", apiPath, strings.ToLower(target)))
			b.WriteString(fmt.Sprintf("    <div id=\"%s-list\">Loading...</div>\n", strings.ToLower(target)))
			b.WriteString("  </div>\n")
		case "create", "input":
			target := a.Target
			if target == "" {
				target = "item"
			}
			apiPath := "/api/" + strings.ToLower(target) + "s"
			b.WriteString(fmt.Sprintf("  <form hx-post=\"%s\" hx-target=\"#%s-list\" hx-swap=\"beforeend\">\n", apiPath, strings.ToLower(target)))
			// Generate input fields from model.
			for _, m := range app.Data {
				if strings.EqualFold(m.Name, target) {
					for _, f := range m.Fields {
						inputType := "text"
						switch f.Type {
						case "number":
							inputType = "number"
						case "email":
							inputType = "email"
						case "date":
							inputType = "date"
						case "boolean":
							inputType = "checkbox"
						}
						required := ""
						if f.Required {
							required = " required"
						}
						b.WriteString(fmt.Sprintf("    <label>%s<input type=\"%s\" name=\"%s\"%s></label>\n", f.Name, inputType, f.Name, required))
					}
				}
			}
			b.WriteString("    <button type=\"submit\">Create</button>\n")
			b.WriteString("  </form>\n")
		case "navigate":
			b.WriteString(fmt.Sprintf("  <a hx-get=\"%s\" hx-target=\"#content\" hx-push-url=\"true\">%s</a>\n", "/"+strings.ToLower(a.Target), a.Text))
		case "delete":
			b.WriteString(fmt.Sprintf("  <button hx-delete=\"/api/%ss/{{.ID}}\" hx-confirm=\"Are you sure?\" hx-target=\"closest .item\" hx-swap=\"outerHTML\">Delete</button>\n", strings.ToLower(a.Target)))
		default:
			b.WriteString(fmt.Sprintf("  <!-- %s: %s -->\n", a.Type, a.Text))
		}
	}

	b.WriteString("</div>\n{{end}}\n")
	return b.String()
}

func generateAlpineJS(app *application) string {
	var b strings.Builder
	b.WriteString("// Alpine.js components for " + app.Name + "\n")
	b.WriteString("document.addEventListener('alpine:init', () => {\n")

	for _, p := range app.Pages {
		name := strings.ToLower(p.Name)
		b.WriteString(fmt.Sprintf("  Alpine.data('%s', () => ({\n", name))
		b.WriteString("    loading: false,\n")
		b.WriteString("    init() {\n")
		b.WriteString("      // Initialize page state\n")
		b.WriteString("    }\n")
		b.WriteString("  }))\n")
	}

	b.WriteString("})\n")
	return b.String()
}

func generateStyles(app *application) string {
	var b strings.Builder
	b.WriteString("/* Generated styles for " + app.Name + " */\n\n")

	// Use theme colors if available.
	b.WriteString(":root {\n")
	if app.Theme != nil && len(app.Theme.Colors) > 0 {
		for name, value := range app.Theme.Colors {
			b.WriteString(fmt.Sprintf("  --%s: %s;\n", name, value))
		}
	} else {
		b.WriteString("  --primary: #3b82f6;\n")
		b.WriteString("  --background: #ffffff;\n")
		b.WriteString("  --text: #1f2937;\n")
	}
	b.WriteString("}\n\n")

	b.WriteString("body { font-family: system-ui, sans-serif; color: var(--text); background: var(--background); margin: 0; padding: 0; }\n")
	b.WriteString("nav { display: flex; gap: 1rem; padding: 1rem; background: var(--primary); }\n")
	b.WriteString("nav a { color: white; text-decoration: none; }\n")
	b.WriteString("main { max-width: 1200px; margin: 0 auto; padding: 2rem; }\n")
	b.WriteString("form { display: flex; flex-direction: column; gap: 0.5rem; max-width: 400px; }\n")
	b.WriteString("label { display: flex; flex-direction: column; gap: 0.25rem; }\n")
	b.WriteString("input { padding: 0.5rem; border: 1px solid #d1d5db; border-radius: 0.25rem; }\n")
	b.WriteString("button { padding: 0.5rem 1rem; background: var(--primary); color: white; border: none; border-radius: 0.25rem; cursor: pointer; }\n")
	b.WriteString(".htmx-indicator { display: none; }\n")
	b.WriteString(".htmx-request .htmx-indicator { display: inline; }\n")

	return b.String()
}

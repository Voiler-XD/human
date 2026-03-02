# human-plugin-htmx

An HTMX + Alpine.js frontend generator plugin for the Human compiler.

## What It Generates

- HTML templates using Go's `html/template` syntax with HTMX attributes
- Navigation with `hx-get` and `hx-push-url` for SPA-like behavior
- Forms with `hx-post` for AJAX submissions
- Display lists with `hx-get` and `hx-trigger="load"` for dynamic content
- Alpine.js components for client-side interactivity
- CSS styles with design token support from the Human theme

## Category

frontend

## Quick Start

```bash
# Build the plugin
make build

# Install into Human
make install

# Verify installation
human plugin list

# Build a .human project — htmx output will appear alongside other generators
human build myapp.human
```

## Protocol

This plugin communicates with the Human compiler via two subcommands:

- `htmx meta` — prints plugin metadata as JSON
- `htmx generate --ir <path> --output <dir>` — generates HTMX templates from IR

## IR Fields Used

| Field | Usage |
|-------|-------|
| `app.Name` | Page title, asset references |
| `app.Pages` | HTML template per page, navigation links |
| `app.Pages[].Content` | HTMX attributes based on action type |
| `app.Data` | Form field generation for create/input actions |
| `app.APIs` | Endpoint paths for hx-get/hx-post |
| `app.Theme` | CSS custom properties from design tokens |

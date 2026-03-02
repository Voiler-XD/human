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
	Name string    `json:"name"`
	Data []*model  `json:"data,omitempty"`
	APIs []*endpoint `json:"apis,omitempty"`
	Auth *auth     `json:"auth,omitempty"`
}

type model struct {
	Name      string      `json:"name"`
	Fields    []*field    `json:"fields,omitempty"`
	Relations []*relation `json:"relations,omitempty"`
}

type field struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Required bool     `json:"required"`
	Unique   bool     `json:"unique"`
	Enum     []string `json:"enum,omitempty"`
}

type relation struct {
	Type   string `json:"type"`
	Target string `json:"target"`
	Field  string `json:"field"`
}

type endpoint struct {
	Name   string `json:"name"`
	Method string `json:"method"`
	Path   string `json:"path"`
	Auth   bool   `json:"auth"`
}

type auth struct {
	Methods []*authMethod `json:"methods,omitempty"`
}

type authMethod struct {
	Type     string `json:"type"`
	Provider string `json:"provider,omitempty"`
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
	resolversDir := filepath.Join(outputDir, "resolvers")
	dirs := []string{outputDir, resolversDir}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return err
		}
	}

	files := map[string]string{
		filepath.Join(outputDir, "schema.graphql"):              generateSchema(app),
		filepath.Join(resolversDir, "index.ts"):                 generateResolverIndex(app),
	}

	// Generate a resolver file per model.
	for _, m := range app.Data {
		filename := strings.ToLower(m.Name) + ".ts"
		files[filepath.Join(resolversDir, filename)] = generateModelResolver(m)
	}

	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

func generateSchema(app *application) string {
	var b strings.Builder
	b.WriteString("# Generated GraphQL schema for " + app.Name + "\n\n")

	// Generate type definitions from data models.
	for _, m := range app.Data {
		b.WriteString(fmt.Sprintf("type %s {\n", m.Name))
		b.WriteString("  id: ID!\n")
		for _, f := range m.Fields {
			gqlType := mapType(f.Type, f.Required)
			if len(f.Enum) > 0 {
				gqlType = m.Name + capitalize(f.Name) + "Enum"
				if f.Required {
					gqlType += "!"
				}
			}
			b.WriteString(fmt.Sprintf("  %s: %s\n", f.Name, gqlType))
		}
		// Add relation fields.
		for _, r := range m.Relations {
			switch r.Type {
			case "has_many":
				b.WriteString(fmt.Sprintf("  %s: [%s!]!\n", strings.ToLower(r.Target)+"s", r.Target))
			case "has_one", "belongs_to":
				b.WriteString(fmt.Sprintf("  %s: %s\n", strings.ToLower(r.Target), r.Target))
			}
		}
		b.WriteString("  createdAt: DateTime!\n")
		b.WriteString("  updatedAt: DateTime!\n")
		b.WriteString("}\n\n")

		// Generate enums.
		for _, f := range m.Fields {
			if len(f.Enum) > 0 {
				enumName := m.Name + capitalize(f.Name) + "Enum"
				b.WriteString(fmt.Sprintf("enum %s {\n", enumName))
				for _, v := range f.Enum {
					b.WriteString(fmt.Sprintf("  %s\n", strings.ToUpper(v)))
				}
				b.WriteString("}\n\n")
			}
		}

		// Generate input types.
		b.WriteString(fmt.Sprintf("input Create%sInput {\n", m.Name))
		for _, f := range m.Fields {
			gqlType := mapType(f.Type, f.Required)
			if len(f.Enum) > 0 {
				gqlType = m.Name + capitalize(f.Name) + "Enum"
				if f.Required {
					gqlType += "!"
				}
			}
			b.WriteString(fmt.Sprintf("  %s: %s\n", f.Name, gqlType))
		}
		b.WriteString("}\n\n")

		b.WriteString(fmt.Sprintf("input Update%sInput {\n", m.Name))
		for _, f := range m.Fields {
			gqlType := mapType(f.Type, false) // all optional for updates
			if len(f.Enum) > 0 {
				gqlType = m.Name + capitalize(f.Name) + "Enum"
			}
			b.WriteString(fmt.Sprintf("  %s: %s\n", f.Name, gqlType))
		}
		b.WriteString("}\n\n")
	}

	// Generate Query type.
	b.WriteString("type Query {\n")
	for _, m := range app.Data {
		lower := strings.ToLower(m.Name)
		b.WriteString(fmt.Sprintf("  %s(id: ID!): %s\n", lower, m.Name))
		b.WriteString(fmt.Sprintf("  %ss(limit: Int, offset: Int): [%s!]!\n", lower, m.Name))
	}
	b.WriteString("}\n\n")

	// Generate Mutation type.
	b.WriteString("type Mutation {\n")
	for _, m := range app.Data {
		lower := strings.ToLower(m.Name)
		b.WriteString(fmt.Sprintf("  create%s(input: Create%sInput!): %s!\n", m.Name, m.Name, m.Name))
		b.WriteString(fmt.Sprintf("  update%s(id: ID!, input: Update%sInput!): %s!\n", m.Name, m.Name, m.Name))
		b.WriteString(fmt.Sprintf("  delete%s(id: ID!): Boolean!\n", lower))
	}
	b.WriteString("}\n\n")

	// Scalar types.
	b.WriteString("scalar DateTime\n")

	return b.String()
}

func generateResolverIndex(app *application) string {
	var b strings.Builder
	b.WriteString("// Generated resolver index for " + app.Name + "\n\n")

	for _, m := range app.Data {
		lower := strings.ToLower(m.Name)
		b.WriteString(fmt.Sprintf("import { %sResolvers } from './%s';\n", lower, lower))
	}

	b.WriteString("\nexport const resolvers = {\n")
	b.WriteString("  Query: {\n")
	for _, m := range app.Data {
		lower := strings.ToLower(m.Name)
		b.WriteString(fmt.Sprintf("    ...%sResolvers.Query,\n", lower))
	}
	b.WriteString("  },\n")
	b.WriteString("  Mutation: {\n")
	for _, m := range app.Data {
		lower := strings.ToLower(m.Name)
		b.WriteString(fmt.Sprintf("    ...%sResolvers.Mutation,\n", lower))
	}
	b.WriteString("  },\n")
	b.WriteString("};\n")

	return b.String()
}

func generateModelResolver(m *model) string {
	name := m.Name
	lower := strings.ToLower(name)

	var b strings.Builder
	b.WriteString(fmt.Sprintf("// Generated resolvers for %s\n\n", name))
	b.WriteString(fmt.Sprintf("export const %sResolvers = {\n", lower))

	// Query resolvers.
	b.WriteString("  Query: {\n")
	b.WriteString(fmt.Sprintf("    %s: async (_: any, { id }: { id: string }, context: any) => {\n", lower))
	b.WriteString(fmt.Sprintf("      // TODO: Fetch %s by ID from database\n", lower))
	b.WriteString(fmt.Sprintf("      return context.db.%s.findUnique({ where: { id } });\n", lower))
	b.WriteString("    },\n")
	b.WriteString(fmt.Sprintf("    %ss: async (_: any, { limit, offset }: { limit?: number; offset?: number }, context: any) => {\n", lower))
	b.WriteString(fmt.Sprintf("      return context.db.%s.findMany({\n", lower))
	b.WriteString("        take: limit ?? 20,\n")
	b.WriteString("        skip: offset ?? 0,\n")
	b.WriteString("      });\n")
	b.WriteString("    },\n")
	b.WriteString("  },\n")

	// Mutation resolvers.
	b.WriteString("  Mutation: {\n")
	b.WriteString(fmt.Sprintf("    create%s: async (_: any, { input }: { input: any }, context: any) => {\n", name))
	b.WriteString(fmt.Sprintf("      return context.db.%s.create({ data: input });\n", lower))
	b.WriteString("    },\n")
	b.WriteString(fmt.Sprintf("    update%s: async (_: any, { id, input }: { id: string; input: any }, context: any) => {\n", name))
	b.WriteString(fmt.Sprintf("      return context.db.%s.update({ where: { id }, data: input });\n", lower))
	b.WriteString("    },\n")
	b.WriteString(fmt.Sprintf("    delete%s: async (_: any, { id }: { id: string }, context: any) => {\n", lower))
	b.WriteString(fmt.Sprintf("      await context.db.%s.delete({ where: { id } });\n", lower))
	b.WriteString("      return true;\n")
	b.WriteString("    },\n")
	b.WriteString("  },\n")

	b.WriteString("};\n")
	return b.String()
}

// mapType converts Human IR field types to GraphQL types.
func mapType(irType string, required bool) string {
	var gql string
	switch irType {
	case "text", "string":
		gql = "String"
	case "number", "integer", "int":
		gql = "Int"
	case "float", "decimal":
		gql = "Float"
	case "boolean", "bool":
		gql = "Boolean"
	case "email":
		gql = "String"
	case "date", "datetime":
		gql = "DateTime"
	case "json":
		gql = "String" // JSON as string; use custom scalar in real impl
	default:
		gql = "String"
	}
	if required {
		gql += "!"
	}
	return gql
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

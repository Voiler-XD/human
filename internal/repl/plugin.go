package repl

import (
	"fmt"
	"strings"

	"github.com/barun-bash/human/internal/cli"
	"github.com/barun-bash/human/internal/plugin"
)

// cmdPlugin dispatches /plugin subcommands.
func cmdPlugin(r *REPL, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(r.out, "Usage: /plugin <list|install|remove|create>")
		return
	}

	switch strings.ToLower(args[0]) {
	case "list", "ls":
		replPluginList(r)
	case "install":
		replPluginInstall(r, args[1:])
	case "remove", "uninstall":
		replPluginRemove(r, args[1:])
	case "create":
		replPluginCreate(r, args[1:])
	default:
		fmt.Fprintf(r.errOut, "Unknown plugin subcommand: %s\n", args[0])
		fmt.Fprintln(r.out, "Usage: /plugin <list|install|remove|create>")
	}
}

func replPluginList(r *REPL) {
	manifests, err := plugin.List()
	if err != nil {
		fmt.Fprintln(r.errOut, cli.Error(fmt.Sprintf("Failed to list plugins: %v", err)))
		return
	}

	if len(manifests) == 0 {
		fmt.Fprintln(r.out, cli.Info("No plugins installed."))
		fmt.Fprintln(r.out, cli.Muted("  Install one with: /plugin install <go-module-path>"))
		return
	}

	fmt.Fprintln(r.out)
	fmt.Fprintln(r.out, cli.Heading("Installed Plugins"))
	fmt.Fprintf(r.out, "  %-20s %-12s %s\n", "NAME", "VERSION", "CATEGORY")
	fmt.Fprintf(r.out, "  %-20s %-12s %s\n", "────", "───────", "────────")
	for _, m := range manifests {
		fmt.Fprintf(r.out, "  %-20s %-12s %s\n", m.Name, m.Version, m.Category)
	}
	fmt.Fprintln(r.out)
}

func replPluginInstall(r *REPL, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(r.errOut, "Usage: /plugin install <go-module-path>")
		fmt.Fprintln(r.errOut, "       /plugin install --binary <path>")
		return
	}

	if args[0] == "--binary" {
		if len(args) < 2 {
			fmt.Fprintln(r.errOut, "Usage: /plugin install --binary <path>")
			return
		}
		fmt.Fprintf(r.out, "  %s Installing plugin from binary: %s\n", cli.Accent("▶"), args[1])
		if err := plugin.InstallFromBinary(args[1]); err != nil {
			fmt.Fprintln(r.errOut, cli.Error(err.Error()))
			return
		}
		fmt.Fprintln(r.out, cli.Success("Plugin installed successfully."))
		return
	}

	source := args[0]
	fmt.Fprintf(r.out, "  %s Installing plugin: %s\n", cli.Accent("▶"), source)
	if err := plugin.Install(source); err != nil {
		fmt.Fprintln(r.errOut, cli.Error(err.Error()))
		return
	}
	fmt.Fprintln(r.out, cli.Success("Plugin installed successfully."))
}

func replPluginRemove(r *REPL, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(r.errOut, "Usage: /plugin remove <name>")
		return
	}

	name := args[0]
	if err := plugin.Uninstall(name); err != nil {
		fmt.Fprintln(r.errOut, cli.Error(err.Error()))
		return
	}
	fmt.Fprintln(r.out, cli.Success(fmt.Sprintf("Plugin %q removed.", name)))
}

func replPluginCreate(r *REPL, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(r.errOut, "Usage: /plugin create <name> [category]")
		return
	}

	name := args[0]
	category := "frontend"
	if len(args) > 1 {
		category = args[1]
	}

	outputDir := name
	fmt.Fprintf(r.out, "  %s Scaffolding plugin: %s (%s)\n", cli.Accent("▶"), name, category)
	if err := plugin.Scaffold(name, category, outputDir); err != nil {
		fmt.Fprintln(r.errOut, cli.Error(err.Error()))
		return
	}
	fmt.Fprintln(r.out, cli.Success(fmt.Sprintf("Plugin project created at ./%s/", name)))
	fmt.Fprintln(r.out, cli.Muted(fmt.Sprintf("  cd %s && make build", name)))
}

// completePlugin provides tab completion for /plugin subcommands.
func completePlugin(_ *REPL, args []string, partial string) []string {
	if len(args) == 0 {
		return completeFromList([]string{"list", "install", "remove", "create"}, partial)
	}

	sub := strings.ToLower(args[0])
	if sub == "remove" || sub == "uninstall" {
		// Complete with installed plugin names.
		manifests, err := plugin.List()
		if err != nil {
			return nil
		}
		names := make([]string, len(manifests))
		for i, m := range manifests {
			names[i] = m.Name
		}
		return completeFromList(names, partial)
	}

	if sub == "install" && len(args) == 1 {
		return completeFromList([]string{"--binary"}, partial)
	}

	if sub == "create" && len(args) == 2 {
		return completeFromList([]string{"frontend", "backend", "database", "infra"}, partial)
	}

	return nil
}

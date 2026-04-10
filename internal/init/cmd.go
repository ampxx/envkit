package init

import (
	"fmt"
	"os"
	"strings"
)

// RunInit is the entry point for the `envkit init` CLI command.
// Args format: [--template <name>] [--dir <path>] [--force]
func RunInit(args []string) error {
	opts := Options{
		Dir:      ".",
		Template: "web",
	}

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--template", "-t":
			if i+1 >= len(args) {
				return fmt.Errorf("--template requires a value")
			}
			i++
			opts.Template = args[i]
		case "--dir", "-d":
			if i+1 >= len(args) {
				return fmt.Errorf("--dir requires a value")
			}
			i++
			opts.Dir = args[i]
		case "--force", "-f":
			opts.Force = true
		case "--list":
			fmt.Println("Available templates:")
			for _, name := range ListTemplates() {
				tmpl := builtinTemplates[name]
				fmt.Printf("  %-12s %s\n", name, tmpl.Description)
			}
			return nil
		default:
			if strings.HasPrefix(args[i], "-") {
				return fmt.Errorf("unknown flag: %s", args[i])
			}
		}
	}

	if err := Scaffold(opts); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return err
	}

	tmpl := builtinTemplates[opts.Template]
	PrintReport(opts, tmpl)
	return nil
}

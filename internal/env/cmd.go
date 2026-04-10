package env

import (
	"fmt"
	"os"
	"text/tabwriter"
)

// RunParse reads a .env file and prints its key-value pairs in a table.
func RunParse(path string) error {
	entries, err := ParseFile(path)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		fmt.Println("No entries found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "KEY\tVALUE")
	fmt.Fprintln(w, "---\t-----")
	for _, e := range entries {
		val := e.Value
		if len(val) > 40 {
			val = val[:37] + "..."
		}
		fmt.Fprintf(w, "%s\t%s\n", e.Key, val)
	}
	return w.Flush()
}

// RunConvert reads a source .env file and writes a normalised copy to dest.
// Existing quoted values are unquoted; values with spaces are re-quoted.
func RunConvert(src, dest string) error {
	entries, err := ParseFile(src)
	if err != nil {
		return fmt.Errorf("convert: read source: %w", err)
	}
	if err := WriteFile(dest, entries); err != nil {
		return fmt.Errorf("convert: write dest: %w", err)
	}
	fmt.Printf("Converted %d entries from %q to %q\n", len(entries), src, dest)
	return nil
}

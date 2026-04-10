package env

import (
	"fmt"
	"sort"
	"strings"
)

// FormatOption controls output behavior for Format.
type FormatOption struct {
	Sorted   bool
	Exported bool // prefix each line with "export "
	Quote    bool // wrap values in double quotes
}

// Format serialises a key→value map back to .env text according to the
// supplied options.  Keys with empty values are included as bare keys.
func Format(vars map[string]string, opt FormatOption) string {
	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	if opt.Sorted {
		sort.Strings(keys)
	}

	var sb strings.Builder
	for _, k := range keys {
		v := vars[k]
		var line string
		switch {
		case opt.Exported && opt.Quote:
			line = fmt.Sprintf("export %s=%q", k, v)
		case opt.Exported:
			line = fmt.Sprintf("export %s=%s", k, v)
		case opt.Quote:
			line = fmt.Sprintf("%s=%q", k, v)
		default:
			line = fmt.Sprintf("%s=%s", k, v)
		}
		sb.WriteString(line)
		sb.WriteByte('\n')
	}
	return sb.String()
}

// Redact returns a copy of vars where every value is replaced with "***".
// Useful for logging or displaying secrets safely.
func Redact(vars map[string]string) map[string]string {
	out := make(map[string]string, len(vars))
	for k := range vars {
		out[k] = "***"
	}
	return out
}

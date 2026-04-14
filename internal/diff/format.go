package diff

import (
	"fmt"
	"io"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorGray   = "\033[90m"
)

// FormatOptions controls output rendering.
type FormatOptions struct {
	Color  bool
	Mask   bool
	ShowUnchanged bool
}

// Render writes a human-readable diff to the given writer.
func Render(w io.Writer, result *Result, opts FormatOptions) {
	fmt.Fprintf(w, "Path: %s  (v%d → v%d)\n", result.Path, result.FromVersion, result.ToVersion)
	fmt.Fprintln(w, strings.Repeat("-", 48))

	for _, c := range result.Changes {
		if c.Type == Unchanged && !opts.ShowUnchanged {
			continue
		}

		switch c.Type {
		case Added:
			val := c.NewValue
			if opts.Mask {
				val = MaskValue(val)
			}
			line := fmt.Sprintf("+ %-24s %s", c.Key, val)
			if opts.Color {
				line = colorGreen + line + colorReset
			}
			fmt.Fprintln(w, line)

		case Removed:
			val := c.OldValue
			if opts.Mask {
				val = MaskValue(val)
			}
			line := fmt.Sprintf("- %-24s %s", c.Key, val)
			if opts.Color {
				line = colorRed + line + colorReset
			}
			fmt.Fprintln(w, line)

		case Modified:
			oldVal, newVal := c.OldValue, c.NewValue
			if opts.Mask {
				oldVal = MaskValue(oldVal)
				newVal = MaskValue(newVal)
			}
			line := fmt.Sprintf("~ %-24s %s → %s", c.Key, oldVal, newVal)
			if opts.Color {
				line = colorYellow + line + colorReset
			}
			fmt.Fprintln(w, line)

		case Unchanged:
			line := fmt.Sprintf("  %-24s (unchanged)", c.Key)
			if opts.Color {
				line = colorGray + line + colorReset
			}
			fmt.Fprintln(w, line)
		}
	}
}

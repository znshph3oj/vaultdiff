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
	Color         bool
	Mask          bool
	ShowUnchanged bool
}

// colorize wraps text with the given ANSI color code if color is enabled.
func colorize(text, color string, enabled bool) string {
	if !enabled {
		return text
	}
	return color + text + colorReset
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
			fmt.Fprintln(w, colorize(line, colorGreen, opts.Color))

		case Removed:
			val := c.OldValue
			if opts.Mask {
				val = MaskValue(val)
			}
			line := fmt.Sprintf("- %-24s %s", c.Key, val)
			fmt.Fprintln(w, colorize(line, colorRed, opts.Color))

		case Modified:
			oldVal, newVal := c.OldValue, c.NewValue
			if opts.Mask {
				oldVal = MaskValue(oldVal)
				newVal = MaskValue(newVal)
			}
			line := fmt.Sprintf("~ %-24s %s → %s", c.Key, oldVal, newVal)
			fmt.Fprintln(w, colorize(line, colorYellow, opts.Color))

		case Unchanged:
			line := fmt.Sprintf("  %-24s (unchanged)", c.Key)
			fmt.Fprintln(w, colorize(line, colorGray, opts.Color))
		}
	}
}

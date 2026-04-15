package audit

import (
	"regexp"
	"strings"
)

// RedactOptions controls which fields are redacted in audit log entries.
type RedactOptions struct {
	// Keys is a list of secret key names whose values should be redacted.
	Keys []string
	// Patterns is a list of regex patterns; values matching any pattern are redacted.
	Patterns []string
}

const redactedPlaceholder = "[REDACTED]"

// Redact returns a copy of the entry with sensitive values masked.
func Redact(entry Entry, opts RedactOptions) Entry {
	if entry.Data == nil {
		return entry
	}

	compiled := compilePatterns(opts.Patterns)
	keySet := toLower(opts.Keys)

	redacted := make(map[string]string, len(entry.Data))
	for k, v := range entry.Data {
		if shouldRedact(k, v, keySet, compiled) {
			redacted[k] = redactedPlaceholder
		} else {
			redacted[k] = v
		}
	}

	copy := entry
	copy.Data = redacted
	return copy
}

// RedactAll applies Redact to every entry in the slice.
func RedactAll(entries []Entry, opts RedactOptions) []Entry {
	out := make([]Entry, len(entries))
	for i, e := range entries {
		out[i] = Redact(e, opts)
	}
	return out
}

func shouldRedact(key, value string, keySet map[string]struct{}, patterns []*regexp.Regexp) bool {
	if _, ok := keySet[strings.ToLower(key)]; ok {
		return true
	}
	for _, re := range patterns {
		if re.MatchString(value) {
			return true
		}
	}
	return false
}

func compilePatterns(patterns []string) []*regexp.Regexp {
	var compiled []*regexp.Regexp
	for _, p := range patterns {
		if re, err := regexp.Compile(p); err == nil {
			compiled = append(compiled, re)
		}
	}
	return compiled
}

func toLower(keys []string) map[string]struct{} {
	m := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		m[strings.ToLower(k)] = struct{}{}
	}
	return m
}

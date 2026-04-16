package audit

import (
	"time"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// TokenEventType describes what kind of token change was observed.
type TokenEventType string

const (
	TokenEventPolicyChange  TokenEventType = "policy_change"
	TokenEventTTLChange     TokenEventType = "ttl_change"
	TokenEventRenewable     TokenEventType = "renewable_change"
	TokenEventLookupSuccess TokenEventType = "lookup_success"
	TokenEventLookupFailure TokenEventType = "lookup_failure"
)

// TokenEvent records an auditable event related to a Vault token.
type TokenEvent struct {
	Timestamp time.Time      `json:"timestamp"`
	Accessor  string         `json:"accessor"`
	EventType TokenEventType `json:"event_type"`
	Detail    string         `json:"detail"`
	User      string         `json:"user,omitempty"`
}

// RecordTokenLookup logs a successful or failed token lookup.
func RecordTokenLookup(l *Logger, info *vault.TokenInfo, user string, lookupErr error) {
	ev := TokenEvent{
		Timestamp: time.Now().UTC(),
		User:      user,
	}
	if lookupErr != nil {
		ev.EventType = TokenEventLookupFailure
		ev.Detail = lookupErr.Error()
	} else if info != nil {
		ev.EventType = TokenEventLookupSuccess
		ev.Accessor = info.Accessor
		ev.Detail = "token lookup succeeded"
	}
	_ = l.Record(Entry{
		Timestamp: ev.Timestamp,
		User:      ev.User,
		Path:      "auth/token/lookup-self",
		Operation: string(ev.EventType),
		Data:      map[string]interface{}{"accessor": ev.Accessor, "detail": ev.Detail},
	})
}

// RecordTokenPolicyChange logs a policy change detected in a token diff.
func RecordTokenPolicyChange(l *Logger, accessor, user string, added, removed []string) {
	_ = l.Record(Entry{
		Timestamp: time.Now().UTC(),
		User:      user,
		Path:      "auth/token/" + accessor,
		Operation: string(TokenEventPolicyChange),
		Data: map[string]interface{}{
			"added":   added,
			"removed": removed,
		},
	})
}

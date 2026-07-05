package smtp

import (
	"strings"
	"testing"

	"github.com/bscott/pm-cli/internal/config"
)

// These tests lock the recipient domain-restriction behavior that pm-cli relies
// on as the authoritative outbound-mail control (the EPA / ephraim daemon
// delegates domain enforcement entirely to pm-cli).
//
// Crucially, validation runs on the assembled *Message* (To/CC/BCC), so it
// covers EVERY CLI input path that populates those fields — the `-t` short
// flag, `--to`, `--cc`, `--bcc`, `to=` template fields, and reply-all. A prior
// wrapper bug validated only `--to`/`to=` at the argv layer and missed `-t` and
// cc/bcc; these tests guard against a regression of that class at the layer that
// actually matters.

func newRestrictedClient(domains ...string) *Client {
	cfg := config.DefaultConfig()
	cfg.Bridge.Email = "epa.vcom@proton.me"
	cfg.Bridge.AllowedDomains = domains
	return NewClient(cfg, "pw")
}

func TestValidateRecipients_AllowsConfiguredDomain(t *testing.T) {
	c := newRestrictedClient("vongerichten.com")
	msg := &Message{To: []string{"matze+asap@vongerichten.com"}}
	if err := c.validateRecipients(msg); err != nil {
		t.Fatalf("expected plus-alias on allowed domain to pass, got: %v", err)
	}
}

func TestValidateRecipients_RejectsForeignRecipientPerField(t *testing.T) {
	// Same foreign recipient must be rejected whether it arrives via To, CC, or
	// BCC — the field it lands in must not matter.
	for _, field := range []string{"to", "cc", "bcc"} {
		t.Run(field, func(t *testing.T) {
			c := newRestrictedClient("vongerichten.com")
			msg := &Message{
				To:  []string{"matze+daily@vongerichten.com"},
				CC:  nil,
				BCC: nil,
			}
			switch field {
			case "to":
				msg.To = append(msg.To, "attacker@evil.com")
			case "cc":
				msg.CC = []string{"attacker@evil.com"}
			case "bcc":
				msg.BCC = []string{"attacker@evil.com"}
			}
			err := c.validateRecipients(msg)
			if err == nil {
				t.Fatalf("expected rejection of foreign recipient in %s field", field)
			}
			if !strings.Contains(err.Error(), "evil.com") {
				t.Fatalf("error should name the rejected domain, got: %v", err)
			}
		})
	}
}

func TestValidateRecipients_DomainMatchIsExact(t *testing.T) {
	c := newRestrictedClient("vongerichten.com")
	cases := []struct {
		name    string
		addr    string
		allowed bool
	}{
		{"exact", "user@vongerichten.com", true},
		{"plus-alias", "matze+triage@vongerichten.com", true},
		{"suffix-trick", "user@vongerichten.com.evil.com", false},
		{"prefix-trick", "user@notvongerichten.com", false},
		{"subdomain", "user@mail.vongerichten.com", false},
		{"no-at", "not-an-address", false},
		{"empty", "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			msg := &Message{To: []string{tc.addr}}
			err := c.validateRecipients(msg)
			if tc.allowed && err != nil {
				t.Fatalf("expected %q allowed, got: %v", tc.addr, err)
			}
			if !tc.allowed && err == nil {
				t.Fatalf("expected %q rejected, but it passed", tc.addr)
			}
		})
	}
}

func TestValidateRecipients_MixedFieldsRejectIfAnyForeign(t *testing.T) {
	c := newRestrictedClient("vongerichten.com")
	msg := &Message{
		To:  []string{"matze+weekly@vongerichten.com"},
		CC:  []string{"matze+monthly@vongerichten.com"},
		BCC: []string{"exfil@attacker.example"},
	}
	if err := c.validateRecipients(msg); err == nil {
		t.Fatal("expected rejection when any single recipient (BCC) is foreign")
	}
}

func TestValidateRecipients_EmptyAllowListMeansNoRestriction(t *testing.T) {
	c := newRestrictedClient() // no allowed_domains configured
	msg := &Message{To: []string{"anyone@anywhere.example"}}
	if err := c.validateRecipients(msg); err != nil {
		t.Fatalf("empty allow-list must impose no restriction, got: %v", err)
	}
}

func TestValidateRecipients_CaseInsensitiveDomain(t *testing.T) {
	c := newRestrictedClient("vongerichten.com")
	msg := &Message{To: []string{"Matze+Asap@Vongerichten.COM"}}
	if err := c.validateRecipients(msg); err != nil {
		t.Fatalf("domain match must be case-insensitive, got: %v", err)
	}
}

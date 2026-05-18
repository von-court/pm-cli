package cli

import (
	"testing"

	"github.com/bscott/pm-cli/internal/imap"
)

func TestParseFieldsValid(t *testing.T) {
	fields, err := parseFields("uid,subject,from_address")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fields) != 3 {
		t.Fatalf("got %d fields, want 3", len(fields))
	}
	if fields[0] != "uid" || fields[1] != "subject" || fields[2] != "from_address" {
		t.Errorf("got %v, want [uid subject from_address]", fields)
	}
}

func TestParseFieldsWithSpaces(t *testing.T) {
	fields, err := parseFields(" uid , subject ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fields) != 2 {
		t.Fatalf("got %d fields, want 2", len(fields))
	}
}

func TestParseFieldsInvalid(t *testing.T) {
	_, err := parseFields("uid,bogus")
	if err == nil {
		t.Fatal("expected error for unknown field")
	}
}

func TestParseFieldsEmpty(t *testing.T) {
	_, err := parseFields("")
	if err == nil {
		t.Fatal("expected error for empty fields")
	}
}

func TestFilterMessageFields(t *testing.T) {
	messages := []imap.MessageSummary{
		{
			UID:         1,
			From:        "Alice",
			FromAddress: "alice@example.com",
			Subject:     "Hello",
			To:          []string{"bob@example.com"},
			MessageID:   "<msg1@example.com>",
		},
	}

	fields := []string{"uid", "from_address", "subject"}
	result := filterMessageFields(messages, fields)

	if len(result) != 1 {
		t.Fatalf("got %d rows, want 1", len(result))
	}

	row := result[0]
	if row["uid"] != uint32(1) {
		t.Errorf("uid = %v, want 1", row["uid"])
	}
	if row["from_address"] != "alice@example.com" {
		t.Errorf("from_address = %v, want alice@example.com", row["from_address"])
	}
	if row["subject"] != "Hello" {
		t.Errorf("subject = %v, want Hello", row["subject"])
	}
	if _, ok := row["from"]; ok {
		t.Error("from should not be present when not requested")
	}
}

func TestFilterMessageFieldsAllFields(t *testing.T) {
	messages := []imap.MessageSummary{
		{
			UID:         42,
			SeqNum:      10,
			From:        "Test User",
			FromAddress: "test@example.com",
			To:          []string{"a@b.com"},
			MessageID:   "<id@example.com>",
			Subject:     "Test",
			Date:        "2026-01-01 00:00",
			DateISO:     "2026-01-01T00:00:00Z",
			Seen:        true,
			Flagged:     false,
		},
	}

	allFields := []string{
		"uid", "seq", "from", "from_address", "to", "message_id",
		"subject", "date", "date_iso", "seen", "flagged",
	}
	result := filterMessageFields(messages, allFields)
	row := result[0]

	if len(row) != len(allFields) {
		t.Errorf("got %d fields, want %d", len(row), len(allFields))
	}
}

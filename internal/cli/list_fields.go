package cli

import (
	"fmt"
	"strings"

	"github.com/bscott/pm-cli/internal/imap"
)

var fieldAccessors = map[string]func(imap.MessageSummary) interface{}{
	"uid":          func(m imap.MessageSummary) interface{} { return m.UID },
	"seq":          func(m imap.MessageSummary) interface{} { return m.SeqNum },
	"from":         func(m imap.MessageSummary) interface{} { return m.From },
	"from_address": func(m imap.MessageSummary) interface{} { return m.FromAddress },
	"to":           func(m imap.MessageSummary) interface{} { return m.To },
	"message_id":   func(m imap.MessageSummary) interface{} { return m.MessageID },
	"in_reply_to":  func(m imap.MessageSummary) interface{} { return m.InReplyTo },
	"subject":      func(m imap.MessageSummary) interface{} { return m.Subject },
	"date":         func(m imap.MessageSummary) interface{} { return m.Date },
	"date_iso":     func(m imap.MessageSummary) interface{} { return m.DateISO },
	"seen":         func(m imap.MessageSummary) interface{} { return m.Seen },
	"flagged":      func(m imap.MessageSummary) interface{} { return m.Flagged },
}

func parseFields(fieldsStr string) ([]string, error) {
	var fields []string
	for _, f := range strings.Split(fieldsStr, ",") {
		f = strings.TrimSpace(f)
		if f == "" {
			continue
		}
		if _, ok := fieldAccessors[f]; !ok {
			return nil, fmt.Errorf("unknown field %q (valid: %s)", f, validFieldNames())
		}
		fields = append(fields, f)
	}
	if len(fields) == 0 {
		return nil, fmt.Errorf("no valid fields specified")
	}
	return fields, nil
}

func filterMessageFields(messages []imap.MessageSummary, fields []string) []map[string]interface{} {
	result := make([]map[string]interface{}, len(messages))
	for i, msg := range messages {
		row := make(map[string]interface{}, len(fields))
		for _, f := range fields {
			row[f] = fieldAccessors[f](msg)
		}
		result[i] = row
	}
	return result
}

func validFieldNames() string {
	names := []string{
		"uid", "seq", "from", "from_address", "to", "message_id",
		"in_reply_to", "subject", "date", "date_iso", "seen", "flagged",
	}
	return strings.Join(names, ", ")
}

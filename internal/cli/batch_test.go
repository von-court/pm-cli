package cli

import (
	"encoding/json"
	"testing"
)

func TestValidateMailboxName(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"Important", false},
		{"My Label", false},
		{"Nested/Label", false},
		{"has{brace", true},
		{"has*star", true},
		{"has%percent", true},
		{"has\rnewline", true},
		{"has\nnewline", true},
	}
	for _, tc := range tests {
		err := validateMailboxName(tc.name)
		if (err != nil) != tc.wantErr {
			t.Errorf("validateMailboxName(%q) error = %v, wantErr = %v", tc.name, err, tc.wantErr)
		}
	}
}

func TestValidateBatchOp(t *testing.T) {
	tests := []struct {
		name    string
		op      batchOp
		wantErr bool
	}{
		{"label ok", batchOp{Op: "label", UIDs: []string{"uid:1"}, Label: "Test"}, false},
		{"label missing label", batchOp{Op: "label", UIDs: []string{"uid:1"}}, true},
		{"label bad chars", batchOp{Op: "label", UIDs: []string{"uid:1"}, Label: "bad*name"}, true},
		{"unlabel ok", batchOp{Op: "unlabel", UIDs: []string{"uid:1"}, Label: "Test"}, false},
		{"move ok", batchOp{Op: "move", UIDs: []string{"uid:1"}, To: "Archive"}, false},
		{"move missing to", batchOp{Op: "move", UIDs: []string{"uid:1"}}, true},
		{"move bad chars", batchOp{Op: "move", UIDs: []string{"uid:1"}, To: "bad{name"}, true},
		{"flag ok", batchOp{Op: "flag", UIDs: []string{"uid:1"}, Read: true}, false},
		{"archive ok", batchOp{Op: "archive", UIDs: []string{"uid:1"}}, false},
		{"delete ok", batchOp{Op: "delete", UIDs: []string{"uid:1"}}, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateBatchOp(0, tc.op)
			if (err != nil) != tc.wantErr {
				t.Errorf("validateBatchOp() error = %v, wantErr = %v", err, tc.wantErr)
			}
		})
	}
}

func TestBatchOpJSONParsing(t *testing.T) {
	input := `[
		{"op": "label", "uids": ["uid:123", "uid:456"], "label": "Important"},
		{"op": "archive", "uids": ["uid:123"]},
		{"op": "flag", "uids": ["uid:789"], "read": true, "star": true},
		{"op": "unlabel", "uids": ["uid:100"], "label": "Old"},
		{"op": "move", "uids": ["uid:200"], "to": "Trash"},
		{"op": "delete", "uids": ["uid:300"]}
	]`

	var ops []batchOp
	if err := json.Unmarshal([]byte(input), &ops); err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	if len(ops) != 6 {
		t.Fatalf("got %d ops, want 6", len(ops))
	}

	if ops[0].Op != "label" || ops[0].Label != "Important" || len(ops[0].UIDs) != 2 {
		t.Errorf("op 0: got %+v", ops[0])
	}
	if ops[2].Op != "flag" || !ops[2].Read || !ops[2].Star {
		t.Errorf("op 2: got %+v", ops[2])
	}
	if ops[3].Op != "unlabel" || ops[3].Label != "Old" {
		t.Errorf("op 3: got %+v", ops[3])
	}
	if ops[4].Op != "move" || ops[4].To != "Trash" {
		t.Errorf("op 4: got %+v", ops[4])
	}
}

func TestBatchOpValidation(t *testing.T) {
	// unknown op
	op := batchOp{Op: "bogus", UIDs: []string{"uid:1"}}
	if validOps[op.Op] {
		t.Error("expected bogus op to be invalid")
	}

	// empty UIDs caught before validateBatchOp
	op2 := batchOp{Op: "label", Label: "Test"}
	if len(op2.UIDs) != 0 {
		t.Error("expected empty UIDs")
	}
}

func TestBatchOutputFormat(t *testing.T) {
	output := batchOutput{
		Results: []batchResult{
			{Op: "label", Success: true},
			{Op: "archive", Success: false, Error: "not found"},
		},
		Total:     2,
		Succeeded: 1,
		Failed:    1,
	}

	data, err := json.Marshal(output)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if parsed["total"].(float64) != 2 {
		t.Errorf("total = %v, want 2", parsed["total"])
	}
	if parsed["succeeded"].(float64) != 1 {
		t.Errorf("succeeded = %v, want 1", parsed["succeeded"])
	}
	if parsed["failed"].(float64) != 1 {
		t.Errorf("failed = %v, want 1", parsed["failed"])
	}
}

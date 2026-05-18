# mail batch — JSON Format

The `mail batch` command reads a JSON array of operations from stdin (or `--file`) and executes them in a single IMAP connection.

## Usage

```bash
echo '[{"op": "label", "uids": ["uid:123"], "label": "Important"}]' | pm-cli mail batch --json
pm-cli mail batch --json --file ops.json
pm-cli mail batch --json --stop-on-error < ops.json
```

## Input Format

A JSON array of operation objects:

```json
[
  {"op": "label",   "uids": ["uid:123", "uid:456"], "label": "Important"},
  {"op": "unlabel", "uids": ["uid:123"],             "label": "Old"},
  {"op": "archive", "uids": ["uid:789"]},
  {"op": "move",    "uids": ["uid:100"],             "to": "Trash"},
  {"op": "flag",    "uids": ["uid:200"],             "read": true, "star": true},
  {"op": "delete",  "uids": ["uid:300"]}
]
```

## Operations

| Op | Required Fields | Description |
|----|----------------|-------------|
| `label` | `uids`, `label` | Add label (copies to Labels/\<name\>) |
| `unlabel` | `uids`, `label` | Remove label (deletes from Labels/\<name\>) |
| `archive` | `uids` | Move to Archive |
| `move` | `uids`, `to` | Move to named mailbox |
| `flag` | `uids`, plus one or more of: `read`, `unread`, `star`, `unstar` | Set/clear message flags |
| `delete` | `uids` | Move to Trash |

## Common Fields

| Field | Type | Description |
|-------|------|-------------|
| `op` | string | Operation name (required) |
| `uids` | string[] | Message selectors, e.g. `"uid:123"` or `"42"` (required) |
| `mailbox` | string | Source mailbox (default: `"INBOX"`) |

## Output Format

```json
{
  "results": [
    {"op": "label", "success": true},
    {"op": "archive", "success": false, "error": "mailbox not found"}
  ],
  "total": 2,
  "succeeded": 1,
  "failed": 1
}
```

## Flags

| Flag | Description |
|------|-------------|
| `--file, -f` | Read operations from file instead of stdin |
| `--stop-on-error` | Stop executing after the first failed operation |

## Limits

- Maximum input size: 10 MB
- Label and mailbox names must not contain IMAP special characters (`{`, `*`, `%`)

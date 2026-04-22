---
name: config-schema-sync
description: Use when reviewing, writing, or committing changes that touch model/config.go or any TOML config struct in this repo. Verifies that every added, renamed, or removed TOML field has a matching update in sesh.schema.json so the JSON schema used by TOML editors (via `#:schema` directive) stays in sync with the Go config model. Trigger whenever a diff touches files under model/, configurator/, or any struct tagged with `toml:"..."`.
---

# Config ↔ Schema Sync Check

The `sesh.schema.json` file at the repo root is the public JSON Schema that TOML editors load via `#:schema https://github.com/joshmedeski/sesh/raw/main/sesh.schema.json`. Any config field that exists in Go but is missing from the schema produces "Additional properties are not allowed" errors for users (see issue #367 for precedent).

## When to run this check

Run whenever a change touches any of:

- `model/config.go`
- Any Go struct with `toml:"..."` tags
- Files under `configurator/`

## Procedure

1. **List the TOML fields added/renamed/removed** in the diff:

   ```bash
   git diff --unified=0 -- model/config.go configurator/ | rg 'toml:"[^"]+"'
   ```

2. **For every changed field, verify `sesh.schema.json` matches.** Map the Go struct to its schema location:

   | Go struct | Schema path |
   |---|---|
   | `Config` | top-level `properties` |
   | `DefaultSessionConfig` | `properties.default_session.properties` |
   | `SessionConfig` | `properties.session.items.properties` |
   | `WindowConfig` | `properties.window.items.properties` |
   | `WildcardConfig` | `properties.wildcard.items.properties` |

   Check each changed field:

   ```bash
   rg '"<field_name>"' sesh.schema.json
   ```

3. **Report discrepancies.** For each field present in Go but missing (or stale) in the schema, either:
   - Add it to the schema with an appropriate `type`, `description`, and `default` where applicable, or
   - Flag it in the review output so the author can decide.

4. **Preserve schema conventions** when editing:
   - `additionalProperties: false` on nested objects — keep it.
   - `description` is required and should mirror the field's purpose (not just restate the name).
   - Include `default` when the Go zero-value is meaningful.
   - Use `type: "string" | "boolean" | "integer" | "array"` matching the Go type.

## Quick diff-check one-liner

```bash
# Extract TOML field names from Go, compare against schema
comm -23 \
  <(rg -o 'toml:"([^,"]+)' -r '$1' model/config.go | sort -u) \
  <(rg -o '"([a-z_]+)":\s*\{' -r '$1' sesh.schema.json | sort -u)
```

Any field printed is in Go but not in the schema — investigate.

## Do not

- Do not bump the schema `$id` URL.
- Do not remove `additionalProperties: false` to "fix" a missing field — add the field instead.
- Do not add fields to the schema that are not yet merged in Go; the schema should track `main`.

# Track D — Whitelisted filters

Skill: build `WHERE` / `ORDER BY` from query params without concatenating user input into SQL. Bound values only. Unknown keys dropped.

Time: 5–7 hours. Builds on B + C (you already have paginated `List`).

Same practice module. Stdlib + your notes app. Illum’s `filter` package is the **answer key**, not something to import. First pass: write a small `Build` yourself. Docs / Effective Go ok. Open Illum only after yours compiles or you’ve been stuck ~20 minutes.

## What you are practicing

You already call Illum’s `filter.ParamsFromQuery` + `filter.Build` on tickets. D is why that whitelist exists:

- Param keys map to **allowed** column + operator (never free-form column names from the client)
- Values are always **bind args** (`?`), never string-interpolated into the SQL
- Sort column comes from the whitelist, not from raw `?sort=...`
- Default `ORDER BY` so pages stay stable when no sort is sent
- A test that a malicious-looking param is **ignored**, not executed

## Build this

Keep pagination (`page` / `limit`). Add filters on `GET /notes`.

### Support at least

| Query param | Meaning |
| --- | --- |
| `title_ilike` | title ILIKE `%value%` |
| `title_eq` | title = value (optional but good) |
| `sort_by=created` + `order=ASC\|DESC` | order by `created_at` (whitelist the sort key; do not take a raw column name) |

Drop any other filter keys silently (e.g. `name_eq=1;drop table notes` does nothing harmful).

### Shape (Illum-like, smaller)

```text
internal/utils/filter/          # or internal/filter/
  filter.go                     Build(raw map[string]string, whitelist) → where, binds, orderBy
  filter_test.go                table-driven tests

internal/domain/notes/
  filters.go                    NotesFilterWhitelist map[string]Rule

internal/repository/postgres/
  notes_repo.go                 List(page, limit, filterParams) applies Build result + pagination

handler                         parse query → map[string]string (or pass c.Request.URL.Query trimmed)
                                still does NOT build SQL
```

Suggested types (you can name them differently):

```go
type Rule struct {
    Column   string   // single column
    Operator string   // "eq", "ilike", or "" for sort-only keys
    Columns  []string // optional: multi-column ilike later — skip if you want
}

type Result struct {
    WhereClause string
    BindValues  []any
    OrderBy     string
}
```

### Rules

- Whitelist lives in **domain**, not hardcoded inside the repo.
- Repo calls `filter.Build(params, notes.NotesFilterWhitelist)` — never concatenates `paramKey` or `paramValue` into operators/columns.
- Always have a default order (e.g. `created_at DESC`) when `OrderBy` is empty.
- Pagination still uses `Count` + `Limit` + `Offset` on the **same** filtered query (count must respect filters).
- Handler / service: still no raw SQL. Keep the dependency arrow: handler → service → repo.

### Do not

- Import Illum’s `internal/utils/filter`.
- Support every operator Illum has. `eq`, `ilike`, and sort are enough.
- Start Track E middleware here.

## Answer key in Illum

| File | What to steal the idea of |
| --- | --- |
| `internal/utils/filter/filter.go` | `Build`, whitelist `Rule`, bind values, skip unknown keys |
| `internal/utils/filter/filter_test.go` | table-driven tests — copy the *style* |
| `internal/domain/class/filters.go` (or `topic/filters.go`) | whitelist in domain, `title_ilike`-style keys, sort keys with empty `Operator` |
| Any list repo using `filter.Build` | append `WhereClause` / `BindValues`, default `ORDER BY` |

## Check it

```bash
# baseline page
curl -s 'http://localhost:9090/notes?page=1&limit=5'

# ilike
curl -s 'http://localhost:9090/notes?title_ilike=hello&page=1&limit=10'

# sort
curl -s 'http://localhost:9090/notes?sort_by=created&order=ASC&page=1&limit=10'

# junk key must not 500 and must not change SQL meaningfully
curl -s 'http://localhost:9090/notes?name_eq=1;drop%20table%20notes&page=1&limit=10'
```

And tests:

```bash
go test ./internal/utils/filter/...   # or wherever you put Build
```

At least one case: unknown key → empty/extra-free where; one case: `title_ilike` → `ILIKE` + bind with `%value%`.

## Done when

- You can explain SQL injection vs bind parameters in one sentence.
- Sort column comes from the whitelist, never from the query string directly.
- Default `ORDER BY` exists when no sort param is sent.
- You have a **test**, not only a working endpoint.
- Filtered `total` matches the filtered set (not the unfiltered table count).

## Skip for D

`OpIn`, multi-column search, `NULLS LAST`, full Illum operator set. Read those in Illum after yours works.

When D is done, E is middleware / request context (you already have a request-id middleware — E deepens auth-style stacks). Or F if you want transactions next. Your call.

# Track B — GORM that you own

Skill: models, Create / First / Where, not-found, pagination. You already have a notes API in memory. Keep the Gin handlers. Replace the store.

Time: 6–8 hours. Builds on A.

Do this in the same practice module. GORM docs are allowed. Illum answer-key files only after yours compiles or you have been stuck 20 minutes. No Cursor completions on the first pass.

## What you are practicing

A taught the request path. B is the query itself — the muscle you skip when every ticket already has a repo.

- A GORM model with tags
- `Create`, `First`, `Where`, `Count`, `Limit`, `Offset`
- `errors.Is(err, gorm.ErrRecordNotFound)` → your own sentinel
- Soft-delete: a deleted row must disappear from List
- Page / limit math you can write without opening a helper

## Build this

Same routes. Postgres via GORM.

| Method | Path | Behavior |
| --- | --- | --- |
| `POST` | `/notes` | Persist the note. Return `201`. |
| `GET` | `/notes/:id` | Load by id. Missing → your `ErrNotFound` → `404`. |
| `GET` | `/notes` | List with `?page=` and `?limit=`. Default page 1, limit 10, cap limit at 100. Response includes items + total (and page/limit if you want the Illum shape). |
| soft-delete | — | Soft-delete one row (GORM `DeletedAt` or an explicit delete). Prove `GET /notes` no longer returns it. `GET /notes/:id` for that id is `404`. |

Model fields at minimum: `id`, `title`, `body` (or keep `description`), `created_at`, and whatever soft-delete needs.

### Rules

- Check `gorm.ErrRecordNotFound` with `errors.Is`. Do not compare error strings.
- Repo / store methods use a pointer receiver.
- Offset is `(page - 1) * limit`. Reimplement that yourself — do not import Illum’s pagination package.
- Empty title still fails before you hit the DB (keep A’s validation).
- No Preload graphs, scopes, hooks, or raw dashboard SQL. Later.

## Answer key in Illum

Look at the shape, not the whole file.

| File | What to steal the idea of |
| --- | --- |
| `internal/domain/class/model.go` | struct tags — not the whole class domain |
| `internal/repository/postgres/class_repo.go` | `Create` and `GetByID` first, then `List` around the pagination / count pattern |
| `internal/utils/pagination/pagination.go` | `Offset()` and `NewResult` math — reimplement, do not copy the package |
| `CLAUDE.md` §6 testing note | GORM `Create` on Postgres emits `RETURNING` → in sqlmock it is a **Query**, not Exec. You do not need sqlmock in B; just know this for later. |

## Check it

```bash
# create
curl -s -D - -H 'Content-Type: application/json' -d '{"title":"hello","description":"x"}' http://localhost:9090/notes

# get
curl -s -D - http://localhost:9090/notes/<id>

# list
curl -s -D - 'http://localhost:9090/notes?page=1&limit=5'

# soft-delete one row somehow (handler, or a one-off in main / a small script), then:
curl -s -D - 'http://localhost:9090/notes?page=1&limit=50'
curl -s -D - http://localhost:9090/notes/<deleted-id>
```

List must not include the deleted note. Get-by-id on it must be `404`.

## Done when

- You can write a model, `db.Create`, `First`, `Where`, `Count`, `Limit`, `Offset` without opening docs.
- You map not-found with `errors.Is(..., gorm.ErrRecordNotFound)`.
- You can say what a pointer receiver on the repo/store is for.
- Soft-delete is proven with the curls above, not assumed.

Skip for B: filters/whitelists (that is D), transactions (F), sqlmock (G).

When B is done, C is the same notes app split into domain interface + service + repo wiring — you already have half of that instinct from A’s `notes` package and `errors.Is`.

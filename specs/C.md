# Track C — One vertical slice

Skill: repository interface in the domain, a thin service, constructor injection, handler only maps HTTP. Same notes API. New boundaries.

Time: 6–8 hours. Builds on A + B.

You already have handlers talking to `*notes.Repository` with GORM inside the same package. C is the Illum split: **domain owns the contract, postgres implements it, service decides, handler speaks HTTP.** Do not rewrite Create/List/Delete logic from scratch — move and wire it.

Same rules as before: GORM/Gin docs ok. Illum answer-key files only after yours compiles or you have been stuck 20 minutes. First pass without Cursor generating the slice for you.

## What you are practicing

B taught queries. C teaches why Illum looks the way it does:

- A `Repository` **interface** in the domain (no `*gorm.DB` in the interface)
- A `Service` that validates and calls the interface
- A postgres (or `internal/notes` store) type that **implements** that interface
- Handler depends on `*Service`, not on the concrete repo
- `errors.Is` in the handler; sentinel errors live in the domain
- `domain` / service code does **not** import `gin` or `gorm`

If you can swap the postgres repo for an in-memory fake and the service still compiles, you got it.

## Target shape

Something like this (names can match your taste; the boundaries cannot):

```text
internal/domain/note/     (or keep internal/notes as domain-only)
  model.go                Note, NoteListRow, request types if domain-owned
  errors.go               ErrNotFound (and friends)
  repository.go           interface { GetByID, Create, List, Delete }
  service.go              NewService(repo Repository), business rules

internal/repository/postgres/   (or internal/notes/postgres)
  note_repo.go            struct { db *gorm.DB }, implements note.Repository
                          maps gorm.ErrRecordNotFound → note.ErrNotFound

internal/handlers/
  notes.go                NewNotesHandler(svc *note.Service, ...), HTTP only

main.go                   db → NewRepo → NewService → NewHandler → routes
```

You can keep a flatter layout if you want (`internal/notes` = domain, `internal/notes/postgres` = impl). What matters:

| Layer | May import | Must not |
| --- | --- | --- |
| domain / service | stdlib, your domain types | `gin`, `gorm` |
| postgres repo | `gorm`, domain package | `gin`, HTTP status codes |
| handler | `gin`, service, response helpers | `gorm`, SQL, building queries |
| main / routes | everything for wiring | business rules |

## Build this

Same HTTP behavior as B. No new product features.

1. **Extract the interface** — methods you already have: `GetByID`, `Create`, `List`, `Delete`.
2. **Move GORM out of the domain package** — the struct with `db *gorm.DB` lives in the repo impl. Domain `Note` model: prefer plain fields (or keep `gorm.Model` only on a persistence struct if you want domain pure — either is fine for C; pure is better).
3. **Add `Service`** — `Create` validates title/description (empty → domain error, not a handler `if`). List/Get/Delete mostly delegate. Wrap unexpected errors: `fmt.Errorf("create note: %w", err)`.
4. **Handler** — bind, call service, `errors.Is` → 400/404/500. No `*notes.Repository` field anymore.
5. **Wire in `main`** — `repo := ...New(db); svc := note.NewService(repo); h := handlers.NewNotesHandler(svc, l)`.

### Optional proof (strongly recommended)

Write a tiny in-memory repo in a `_test.go` or `internal/notes/memory` that implements the same interface. Construct `NewService(memRepo)` in a test or a short `main` path. If it compiles and Create/Get work, the interface is real.

## Answer key in Illum

Look at the shape, not the whole file.

| File | What to steal the idea of |
| --- | --- |
| `internal/routes/class_routes.go` | `NewRepo` → `NewService` → `NewHandler` in one place |
| `internal/domain/class/service.go` | `NewService(repo)`, validate then repo — rules not in Gin |
| `internal/domain/class/service.go` `CreateClass` | business rules live here |
| `CLAUDE.md` §2 and §5 | layer rules; sentinel → HTTP with `errors.Is` |

Do not copy TxManager, filters, or institute stacks. That is later.

## Check it

All B curls still pass:

```bash
curl -s -D - -H 'Content-Type: application/json' -d '{"title":"hello","description":"x"}' http://localhost:9090/notes
curl -s -D - http://localhost:9090/notes/<id>
curl -s -D - 'http://localhost:9090/notes?page=1&limit=5'
curl -s -D - -X DELETE http://localhost:9090/notes/<id>
curl -s -D - http://localhost:9090/notes/<id>
```

Plus:

- Empty title still `400`, and the decision happens in the **service** (or domain validate), not as the only check in the handler.
- `grep -R 'gorm.io' internal/domain` (or your domain package) → no hits, except maybe a persistence DTO if you chose that split.
- `grep -R 'gin-gonic' internal/notes` / domain / service → no hits.
- Handler file has no `*gorm.DB` and no `repo.db`.

## Done when

- Domain does not import `gin` or `gorm`.
- Handler never builds a query; repo never returns HTTP status.
- You can add a second fake repo (in-memory) and the service still compiles.
- You wrap errors as `fmt.Errorf("create note: %w", err)` where it helps.
- You can explain out loud: “handler parses and maps status; service decides; repo queries.”

## Skip for C

DTOs for every field, Swagger, institute ID, filters (D), middleware stacks (E), transactions (F).

When C is done you will recognize Illum PRs faster: “this query is in the handler” and “this repo is doing another domain’s job” become obvious. Then D (whitelist filters) or E (middleware) — your call; both sit on top of this slice.

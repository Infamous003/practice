# Track E — Middleware and request context

Skill: compose `gin.HandlerFunc`s, abort the chain, store request-scoped data with a typed key (not a bare string). Fake auth is enough — no Keycloak.

Time: smaller than D/F. About 3–5 hours if you stay focused. Builds on A (you already have request-id).

Same practice notes app. Gin docs ok. Illum answer-key files only after yours compiles or you’ve been stuck ~20 minutes. First pass without Cursor generating the stack.

## What you are practicing

Illum attaches behavior with `g.Use(middleware.SomeStack()...)`. Until you write a stack, “middleware” is a file someone else owns.

- `func Auth() gin.HandlerFunc { return func(c *gin.Context) { ... } }`
- `c.Next()` runs the rest of the chain; `c.Abort()` stops it after you wrote a response
- Unexported / typed context keys (you already use `RequestIDKey` — keep that habit)
- A **stack** = `[]gin.HandlerFunc` spread with `...`
- Handler reads the user from context; it does not re-parse the auth header

## Build this

Keep `/health` public (request-id only). Protect notes.

### Fake auth

Header: `Authorization: Bearer <user-id>`  
Hardcoded users are fine, e.g.:

| Bearer token | role |
| --- | --- |
| `user-1` | `user` |
| `admin-1` | `admin` |

No JWT. No Keycloak. A map in middleware is the point.

### Middleware

1. **`Auth()`** — missing/unknown bearer → `401` + `c.Abort()`. On success, set a small user value on the Gin context (id + role), then `c.Next()`.
2. **`RequireRole(roles ...string)`** — if the user role is not in the list → `403` + `c.Abort()`.
3. **Optional stack helper** — e.g. `AuthenticatedStack()`, `AdminStack()` returning `[]gin.HandlerFunc`.

Also tighten request-id while you’re here (from A leftovers):

- Set response header `X-Request-ID` to the same value as the envelope
- Bad/missing incoming id → mint a new one (do **not** panic on parse failure)

### Routes

| Route | Access |
| --- | --- |
| `GET /health` | no auth |
| `GET /notes`, `GET /notes/:id` | any authenticated user |
| `POST /notes`, `DELETE /notes/:id` | `admin` only |

Wiring in `main` (Illum style):

```go
notes := r.Group("/notes")
notes.Use(middleware.AuthenticatedStack()...)
{
    notes.GET("", h.List...)
    notes.GET("/:id", h.Get...)
}

adminNotes := r.Group("/notes")
adminNotes.Use(middleware.AdminStack()...)
{
    adminNotes.POST("", h.Create...)
    adminNotes.DELETE("/:id", h.Delete...)
}
```

(or one group with per-route `RequireRole` — either is fine)

### Handler rule

Create/Delete may read the user from context for logging. They must **not** call `c.GetHeader("Authorization")` again. Auth already decided.

## Answer key in Illum

| File | What to steal the idea of |
| --- | --- |
| `internal/middleware/request_id.go` | smallest correct middleware: header in, set context, header out, `c.Next()` |
| `internal/middleware/stacks.go` | stacks are slices; spread with `...` |
| `internal/middleware/auth.go` | abort + error response, then stop reading (ignore Keycloak/session) |
| `internal/context/request.go` | typed keys + helpers; why not string keys |

## Check it

```bash
# public
curl -s -D - http://localhost:9090/health

# no auth → 401
curl -s -D - http://localhost:9090/notes

# user can list, cannot create
curl -s -D - -H 'Authorization: Bearer user-1' 'http://localhost:9090/notes?page=1&limit=5'
curl -s -D - -H 'Authorization: Bearer user-1' -H 'Content-Type: application/json' \
  -d '{"title":"x","description":"y"}' http://localhost:9090/notes

# admin can create
curl -s -D - -H 'Authorization: Bearer admin-1' -H 'Content-Type: application/json' \
  -d '{"title":"x","description":"y"}' http://localhost:9090/notes

# request id echoes
curl -s -D - -H 'X-Request-ID: 11111111-1111-1111-1111-111111111111' http://localhost:9090/health
```

Expect: `401` / `403` / `201` in the right places; `X-Request-ID` on the response.

## Done when

- You can write `func Auth() gin.HandlerFunc { return func(c *gin.Context) { ... } }` from memory
- You know `c.Next()` vs `c.Abort()` after writing a response
- You will not use bare string keys in `context.WithValue` / Gin `Set` for this (typed key)
- You can say what `InstituteAdminStack` is without opening Illum: “auth + user-type gate + … as a slice of handlers”

## Skip for E

Real JWT, Keycloak, permission catalogs, institute resolution. A hardcoded user map is the point.

When E is done: **F** (TxManager) is the high-leverage one you already understand conceptually. **D** (filters) can stay parked until you care again.

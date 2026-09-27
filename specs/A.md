# Track A — Gin request path

Skill: write a Gin app from a blank `main.go`. Bind a body, reject a bad one, return a typed envelope, stop the chain when you must.

Time: a couple of hours. No database. No Cursor completions on the first pass.

Do this in `~/code/go-fluency`, not in illum-backend. The files below are the answer key. Open them only after yours compiles, or after you have been stuck for 20 minutes.

## What you are practicing

You already know an HTTP router (Chi). This is the same idea in the dialect this repo uses, written by you:

- `gin.New()`, a group, `GET` / `POST`
- `ShouldBindJSON`, `c.Param`
- a middleware that runs before the handler (`c.Next()`)
- abort versus return
- the handler does not own business rules

## Build this

One module. In-memory map is fine. Nothing persisted.

| Method | Path | Behavior |
| --- | --- | --- |
| `GET` | `/health` | `200`, body says it is up |
| `POST` | `/notes` | bind JSON `{ "title": "...", "body": "..." }`. Empty or missing title → `400`. Otherwise store it and return `201` with the note, including an id you generated. |
| `GET` | `/notes/:id` | `200` with the note, or `404` if that id was never created |

Every JSON response, success or error, uses this envelope:

```json
{ "success": true, "data": {}, "request_id": "..." }
```

```json
{ "success": false, "error_message": "...", "request_id": "..." }
```

`request_id` comes from middleware, not from the handler inventing it inline.

### Middleware

Read `X-Request-ID`. If the header is missing, generate an id. Put that same id on the response header and in the JSON envelope. Then call the rest of the chain.

### Rules

- Empty title is a validation result, not a panic and not a `500`.
- After you write an error response, the rest of the chain must not run. Know whether that means `c.Abort()`, a bare `return`, or both, and why.
- The "is this title allowed?" check does not live in the route registration. The handler parses and maps status. A small function decides whether the title is valid.
- No Swagger, CORS, auth, or Keycloak.

## Answer key

Look at the shape, not the whole file.

| File | What to steal the idea of |
| --- | --- |
| `internal/app/router.go` | `gin` engine, `Group`, `Use()` |
| `internal/middleware/request_id.go` | `HandlerFunc`, read header, set header, `c.Next()` |
| `internal/utils/response/response.go` | `c.JSON` envelope. Do not copy the package. |
| `internal/handlers/user_auth_handler.go` | bind, then map the result to a status. Ignore Keycloak. |

## Check it

```bash
curl -s -D - http://localhost:8080/health
curl -s -D - -H 'Content-Type: application/json' -d '{"title":"","body":"x"}' http://localhost:8080/notes
curl -s -D - -H 'Content-Type: application/json' -H 'X-Request-ID: abc' -d '{"title":"hello","body":"x"}' http://localhost:8080/notes
curl -s -D - http://localhost:8080/notes/<id-from-previous>
curl -s -D - http://localhost:8080/notes/missing
```

The `400` and `201` responses must both carry `request_id`, and when you sent `X-Request-ID: abc` that value comes back.

## Done when

- You can write `gin.New()`, `Group`, `GET` / `POST` from memory.
- You know when to `c.Abort()` versus `return` after writing a response.
- You can use `c.Param`, `c.Query`, and `ShouldBindJSON` without opening docs.
- You can explain why the handler must not contain the business rule.

Next morning: delete the file and rewrite it in 30–40 minutes. Then start B.

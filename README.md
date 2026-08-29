# State Welfare Card Registration (MVP1)

Counter-side registration for Thailand's State Welfare Card. A registrar fills in the application on behalf of the applicant; the system validates format, verifies the address codes are internally consistent, and blocks duplicates before writing.

**MVP1 does not decide eligibility.** It accepts a complete application and prevents duplicates. Eligibility is the downstream system's job.

**Stack** — Go 1.27 + gin + pgx/v5 · PostgreSQL 16 · Vue 3 `<script setup>` + TypeScript + Vite + Pinia + zod + Tailwind v4

## Getting started

Requires Docker, Go 1.27+, Node 20+, and [golang-migrate](https://github.com/golang-migrate/migrate).

```bash
cp .env.example .env
make up && make migrate && make seed
```

Then run the API and the web app in two terminals:

```bash
make api   # http://localhost:8080
```

```bash
make web   # http://localhost:5173
```

Open <http://localhost:5173> and you land on the registration form. There is no login screen — pick a registrar from the dropdown in the top right.

| Command | |
| --- | --- |
| `make up` / `make down` | Start / stop PostgreSQL |
| `make migrate` / `make rollback` | Apply / roll back one migration |
| `make seed` | Load `backend/db/seed/*.sql` |
| `make psql` | psql inside the container |
| `make api` / `make web` | Run API / frontend |
| `make test` | `go test ./...` |

## Configuration

Copy `.env.example` to `.env`. `DATABASE_URL`, `HASH_PEPPER`, and `ENC_KEY` have no defaults — the app refuses to start without them.

| Variable | Default | |
| --- | --- | --- |
| `APP_ENV` | `development` | `production` forbids StubAuth |
| `HTTP_PORT` | `8080` | The Vite dev proxy follows this too |
| `CORS_ORIGIN` | `http://localhost:5173` | |
| `STUB_AUTH_ENABLED` | `false` | |
| `STUB_DEFAULT_USER_ID` | — | User assumed when no header is sent |
| `HASH_PEPPER` | — | Pepper for `national_id_hash` |
| `ENC_KEY` | — | AES-256 key, 64 hex chars |
| `MAX_INFLIGHT` | `512` | Concurrent requests before shedding with `503 SYS002`; `0` disables |
| `MAX_BODY_BYTES` | `1048576` | Max request body |

Generate real secrets with `openssl rand -hex 32`. **Changing `HASH_PEPPER` orphans every existing record** (the lookup hash changes) and **changing `ENC_KEY` makes existing ciphertext undecryptable.** Both need a rotation plan.

To run elsewhere: `HTTP_PORT=9090 make api` and `HTTP_PORT=9090 make web`, or point the frontend at another host with `API_TARGET=http://host:port`.

## StubAuth

No real authentication yet. Identity comes from a header:

```
X-Debug-User-Id: <uuid from the users table>
```

Absent, `STUB_DEFAULT_USER_ID` is used. The user switcher reads `GET /api/v1/dev/users`. Seeded users: `somchai.ktb`, `somying.baac`, `admin`.

The organization and user on an application **always come from the actor resolved by middleware, never from the request body.** Under `APP_ENV=production` the dev route is not registered and config validation fails if StubAuth is enabled.

## API

| Method | Path | |
| --- | --- | --- |
| GET | `/health` | Also pings the database; registered before the in-flight limiter so probes never get shed |
| GET | `/api/v1/ref/{provinces,districts,subdistricts}` | `public, max-age=86400` + ETag |
| GET | `/api/v1/ref/address-search?q=บางจาก` | `private, max-age=60` |
| POST | `/api/v1/applications` | Requires an actor |
| GET | `/api/v1/dev/users` | Only while StubAuth is enabled |

Every response uses one envelope: `data`, `statusCode` (`"0"` success / `"1"` failure — a single digit, not `"0000"`), `statusDescription`, and on failure `errorCode`, `errorMessage`, `fieldErrors`.

```json
"fieldErrors": [
  { "field": "personal.national_id", "code": "VAL001", "message": "..." },
  { "field": "financial.assets[2].unit", "code": "VAL013", "message": "..." }
]
```

`field` is a dot path that **matches the request payload exactly** — the frontend uses it to highlight the offending input, so renaming a field on one side requires renaming it on the other.

Codes: `VAL000`–`VAL015` format, `DUP001` duplicate, `AUTH001`, `SYS001` internal, `SYS002` overloaded (sent with `Retry-After`). Full list in `backend/pkg/apperror/apperror.go`.

## How a submission is validated

```
L1 format  →  L2 own database  →  L3 external system  →  L4 write
```

A failing layer stops the pipeline, but **within a layer every error is collected before returning.** Never bail on the first problem: the applicant is standing at the counter, and one field per round trip is cruelty dressed up as validation.

- **L2** verifies the province/district/subdistrict/postal codes actually belong together rather than trusting the browser, and runs the duplicate check **before** any external call (those cost money per request).
- **L3** is a stub. If the identity provider is down the application is still `SUBMITTED` and still returns 201 — only `verification_method` becomes `PENDING_VERIFICATION`. Someone else's outage is not the applicant's fault.
- **L4** is one transaction: upsert citizen, then insert application, household members, income, assets, liabilities, verification, status history.

## National ID data

- `national_id_hash` — HMAC-SHA256 with a pepper. Used for lookup and duplicate detection, protected by UNIQUE.
- `national_id_enc` — AES-256-GCM stored as `v1|nonce|ciphertext`, with the hash as AAD.
- `laser_id` is **validated and thrown away. It is never persisted.**

Lifetime cards and unreadable laser codes are handled by picking a verification reason instead, with free text (min 10 characters) behind an "other" option.

## Duplicate prevention

Two independent layers:

1. `UNIQUE` on `citizens.national_id_hash`
2. A partial unique index on `(citizen_id, fiscal_year)` covering applications still in play

A duplicate returns `409 DUP001` with the registering organization, date, application number, and appeal deadline. Concurrent submissions of the same ID resolve to exactly one row — one request gets 201, the rest get 409.

The fiscal year is derived from the current date on the frontend and rolls over on 1 October. Registrars cannot pick it, because the wrong year silently defeats the duplicate check.

## Things that will bite you

- **Roll back the transaction before querying for conflict details** when building a 409, or the connection is still aborted and you get a 500.
- **Money is never a float.** `amount` arrives as a `json.Number` string, capped at 12 digits plus 2 decimals; over that is `400 VAL011`, not a 500 at INSERT.
- **`unit` is bound one-to-one to `asset_type`** — see the `assetUnit` map in `internal/domain/application.go`. A mismatch is `VAL013`.
- **The domain layer imports no framework.** No gin, no pgx under `internal/domain/`.
- **Adding a form field means touching four places:** `types/form.ts` → `emptyForm()` → **the zod shape in `schemas/registration.ts`** → `buildPayload()`. Forget the zod shape and the field is stripped before `superRefine` sees it: validation silently does nothing, with no error to trace.
- **`address-search` must not use `public, max-age=86400`** like the other `/ref/*` endpoints. Browsers cache an empty result for a day and the search looks broken while the data sits right there.

## Layout

```
backend/
  cmd/api/        entry point, HTTP timeouts, graceful shutdown
  internal/
    config/       env loading and validation
    domain/       business rules only — appno, birthdate, laserid, nationalid
    dto/          request and response shapes
    handler/      routes and gin handlers
    middleware/   auth, cors, cache, recovery, in-flight and body limits
    repository/   pgx queries
    service/      the four validation layers
    verifier/     external identity verification (stub)
  pkg/            apperror (codes), httpx (envelope), idcrypto (AES-256-GCM)
  db/             migrations, seed

frontend/src/
  api/            axios client (30s timeout) and endpoint wrappers
  components/     common/ (buttons, dialogs, date picker, type+amount row)
                  fields/ (FieldText / FieldNumber / FieldSelect)
                  registration/ (form sections)
  composables/    useRegistrationForm, useFormContext
  constants/      dropdown options and form defaults
  schemas/        zod · stores/ pinia · views/ RegisterFormView, SuccessView
```

Address reference data in `backend/db/seed/01_ref_address.sql` covers 77 provinces, 930 districts, and 7,452 subdistricts with postal codes, using official Department of Provincial Administration codes ([source](https://github.com/kongvut/thai-province-data)). It is generated — do not edit by hand.

## Capacity

Measured with a throwaway load generator (not committed) posting unique, checksum-valid national IDs against a CPU-capped Postgres container:

- **~1,900 applications/second** on API 2 vCPU + database 4 vCPU with a simulated 200 ms identity-provider call. Without that call, ~2,900/s.
- **The database is the bottleneck.** Doubling database CPU from 2 to 4 gained 24%; going to 8 gained 0.2%; giving the API more CPU changed nothing.
- **Replicas help only up to a point** — 1→3 improved throughput 38%, 3→6 made it worse. They all queue on one Postgres.
- Throughput stayed flat from 50 to 800 concurrent while latency grew linearly: the system degrades by queueing, never by dropping. 20,000 concurrent writes produced exactly 20,000 rows with zero duplicates.

**Suggested production spec:** API 2 vCPU / 2 GB × 3 replicas, PostgreSQL 4 vCPU / 8 GB. Do not go below 2 vCPU on the database — at 1 vCPU throughput halves and p99 triples.

**Watch the connection math.** Postgres defaults to `max_connections = 100` and every replica opens its own pool. Six replicas × a 24-connection pool returned 500 for half of all requests; the same load with pools sized to stay under 100 succeeded completely. Keep `maxReplicas × pool size` well under the limit, or put PgBouncer in front.

For a normal registration window none of this is close to binding: 1,000,000 applicants over 30 days averages 0.4 req/s. A truly simultaneous burst is a different problem — draining 1,000,000 at once takes about 6 minutes, longer than any load balancer timeout, and no amount of hardware fixes it because one database serializes the writes. Stagger arrivals instead (by last digit of the national ID, or a virtual waiting room).

**Caveats:** everything ran on one machine over loopback, no TLS, runs lasted seconds rather than hours, and traffic was all-success. **If the real identity provider rate-limits to 50 req/s, then 50 req/s is the true ceiling regardless of any number above.**

## Tests

```bash
make test                          # repository and service tests need `make up`
cd frontend && npx vue-tsc --noEmit # no frontend test suite yet
```

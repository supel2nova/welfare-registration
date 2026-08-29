# State Welfare Card Registration (MVP1)

Counter-side registration system for Thailand's State Welfare Card. A registrar at a bank branch or local government office fills in the application on behalf of the applicant. The system validates the data format, verifies that the address codes are internally consistent, and blocks duplicate applications before writing anything.

**MVP1 does not decide eligibility.** It never judges whether an applicant passes the income or asset thresholds — it only accepts a complete application and prevents duplicates. Eligibility is the downstream system's job.

## Stack

| Layer | Technology |
| --- | --- |
| Backend | Go 1.27, gin, pgx/v5 |
| Database | PostgreSQL 16 (Docker), golang-migrate |
| Frontend | Vue 3 `<script setup>`, TypeScript, Vite, Pinia, vue-router, zod, Tailwind v4 |

## Getting started

Requires Docker, Go 1.27+, Node 20+, and [golang-migrate](https://github.com/golang-migrate/migrate).

```bash
cp .env.example .env
make up        # start PostgreSQL
make migrate   # create all tables
make seed      # load address reference data, organizations, users
```

Then run the API and the web app in two terminals:

```bash
make api   # http://localhost:8080
```

```bash
make web   # http://localhost:5173
```

Open <http://localhost:5173> and you land straight on the registration form. There is no login screen — pick a registrar from the dropdown in the top right instead (see [StubAuth](#stubauth)).

## Commands

| Command | What it does |
| --- | --- |
| `make up` / `make down` | Start / stop PostgreSQL |
| `make migrate` | Apply migrations |
| `make rollback` | Roll back one migration |
| `make seed` | Load every `.sql` file in `backend/db/seed/` |
| `make psql` | Open psql inside the container |
| `make api` | Run the API |
| `make web` | Run the frontend dev server |
| `make test` | `go test ./...` |
| `make tidy` | `go mod tidy` |

## Configuration

Everything lives in `.env` (copy from `.env.example`). Three variables have no default and the app refuses to start without them: `DATABASE_URL`, `HASH_PEPPER`, `ENC_KEY`.

| Variable | Default | Notes |
| --- | --- | --- |
| `APP_ENV` | `development` | Set to `production` and StubAuth can no longer be enabled |
| `HTTP_PORT` | `8080` | The Vite dev proxy reads this too, so `make web` follows it |
| `DATABASE_URL` | — | |
| `CORS_ORIGIN` | `http://localhost:5173` | |
| `STUB_AUTH_ENABLED` | `false` | |
| `STUB_DEFAULT_USER_ID` | — | User assumed when no header is sent |
| `HASH_PEPPER` | — | Pepper for `national_id_hash` |
| `ENC_KEY` | — | AES-256 key, 64 hex characters |

The `HASH_PEPPER` and `ENC_KEY` values in `.env.example` are development throwaways. Generate real ones before any deployment:

```bash
openssl rand -hex 32
```

**Changing `HASH_PEPPER` orphans every existing record** because the lookup hash changes across the whole table, and **changing `ENC_KEY` makes existing ciphertext undecryptable.** Both need a key-rotation plan before you touch them.

### Running on a different port

`HTTP_PORT` is picked up by both the API and the Vite dev proxy, so this is enough:

```bash
HTTP_PORT=9090 make api
HTTP_PORT=9090 make web
```

Or set `HTTP_PORT` in `.env` — the Makefile exports it to both processes. To point the frontend at an API somewhere else entirely, set `API_TARGET=http://host:port` instead.

## StubAuth

MVP1 has no real authentication yet. The registrar's identity comes from a header:

```
X-Debug-User-Id: <uuid from the users table>
```

If the header is absent, `STUB_DEFAULT_USER_ID` is used. The web app has a user switcher in the top right that reads its list from `GET /api/v1/dev/users`.

The organization and user recorded on an application **always come from the actor resolved by middleware, never from the request body.** Even if a client sends them, they are ignored.

When `APP_ENV=production`, the `/api/v1/dev/users` route is not registered at all and StubAuth is off. Config validation fails outright if `STUB_AUTH_ENABLED=true` is combined with production.

Seeded users: `somchai.ktb`, `somying.baac`, `admin`.

## API

| Method | Path | Notes |
| --- | --- | --- |
| GET | `/health` | Also pings the database |
| GET | `/api/v1/ref/provinces` | `public, max-age=86400` + ETag |
| GET | `/api/v1/ref/districts?province_code=10` | same |
| GET | `/api/v1/ref/subdistricts?district_code=1008` | same |
| GET | `/api/v1/ref/address-search?q=บางจาก` | `private, max-age=60` |
| POST | `/api/v1/applications` | Requires an actor |
| GET | `/api/v1/dev/users` | Only while StubAuth is enabled |

Every response uses the same envelope:

```json
{
  "data": null,
  "statusCode": "1",
  "statusDescription": "ข้อมูลไม่ถูกต้อง",
  "errorCode": "VAL001",
  "errorMessage": "ข้อมูลไม่ถูกต้อง",
  "fieldErrors": [
    { "field": "personal.national_id", "code": "VAL001", "message": "..." },
    { "field": "financial.assets[2].unit", "code": "VAL013", "message": "..." }
  ]
}
```

`statusCode` is `"0"` on success and `"1"` on failure — a single digit, not `"0000"`.

`fieldErrors[].field` is a dot path that **matches the request payload exactly.** The frontend uses that path to highlight the offending input directly, so renaming a field on one side requires renaming it on the other.

Error codes: `VAL000`–`VAL015` (format), `DUP001` (duplicate), `AUTH001`, `SYS001`. The full list is in `backend/pkg/apperror/apperror.go`.

## The four validation layers

```
L1 format  →  L2 own database  →  L3 external system  →  L4 write
```

A failing layer stops the pipeline, but **within a single layer every error is collected before returning.** Never bail out on the first problem: the applicant is standing at the counter, and making them fix one field per round trip is cruelty dressed up as validation.

- **L2** verifies that the province, district, subdistrict, and postal codes actually belong together rather than trusting what the browser sent, and runs the duplicate check **before** any external call.
- **L3** is a stub. If the identity provider is down, the application is still `SUBMITTED` and still returns 201 — only `verification_method` becomes `PENDING_VERIFICATION`. An outage on someone else's system is not the applicant's fault.
- **L4** runs in a single transaction: upsert the citizen, then insert the application, household members, income sources, assets, liabilities, and status history.

## Handling national ID data

- `national_id_hash` — HMAC-SHA256 of the ID plus a pepper. Used for lookup and duplicate detection, protected by a UNIQUE constraint.
- `national_id_enc` — AES-256-GCM stored as `v1|nonce|ciphertext`, with `national_id_hash` as the AAD.
- `laser_id` (the code on the back of the card) is **validated and thrown away. It is never persisted,** under any circumstance.

Lifetime cards and cards whose laser code is too damaged to read are handled by selecting a verification reason (with a free-text option requiring at least 10 characters). The application is accepted normally.

## Duplicate prevention

Two independent layers, because one is not enough:

1. `UNIQUE` on `citizens.national_id_hash`
2. A partial unique index on `(citizen_id, fiscal_year)` covering only applications still in play

A duplicate returns `409 DUP001` along with which organization already registered the applicant, when, under what application number, and the appeal deadline.

The fiscal year is derived from the current date on the frontend and rolls over on 1 October. Registrars cannot pick it, because picking the wrong year would silently defeat the duplicate check.

## Things that will bite you

- **Roll back the transaction before querying for conflict details** when building a 409. Otherwise the connection is still in an aborted state and you get a 500 instead.
- **Money is never a float.** The Go side accepts `amount` as a `json.Number` (a string), capped at 12 digits plus 2 decimals. Exceeding that is a `400 VAL011`, not a 500 at INSERT time.
- **`unit` is bound one-to-one to `asset_type`.** See the `assetUnit` map in `backend/internal/domain/application.go`. A mismatched pair is `VAL013`.
- **The domain layer imports no framework.** No gin, no pgx anywhere under `internal/domain/`.
- **Adding a form field means touching four places:** `types/form.ts` → `emptyForm()` → **the zod shape in `schemas/registration.ts`** → `buildPayload()`. Forget the zod shape and the field is stripped before `superRefine` ever sees it, so validation silently does nothing and there is no error to trace.
- **`address-search` must not use `public, max-age=86400`** like the other `/ref/*` endpoints. Browsers will cache an empty result set for a full day and the search will appear broken while the data is sitting right there.

## Project layout

```
backend/
  cmd/api/            entry point
  internal/
    config/           env loading and validation
    domain/           business rules only, no framework imports
      appno/          application number WC-YYYY-nnnnnnn
      birthdate/      partial birth dates and age calculation
      laserid/        laser code format check
      nationalid/     mod-11 checksum
    dto/              request and response shapes
    handler/          routes and gin handlers
    middleware/       auth, cors, cache, recovery
    repository/       pgx queries
    service/          the four validation layers, application assembly
    verifier/         external identity verification (stub)
  pkg/
    apperror/         error codes and fieldErrors
    httpx/            response envelope
    idcrypto/         AES-256-GCM
  db/
    migrations/       golang-migrate
    seed/             addresses, organizations, users

frontend/src/
  api/                axios client and endpoint wrappers
  components/
    common/           buttons, icons, dialogs, date picker, type+amount row
    fields/           FieldText / FieldNumber / FieldSelect
    registration/     the form's sections
  composables/        useRegistrationForm, useFormContext
  constants/          dropdown options and form defaults
  schemas/            zod
  stores/             pinia
  views/              RegisterFormView, SuccessView
```

## Address reference data

`backend/db/seed/01_ref_address.sql` contains 77 provinces, 930 districts, and 7,452 subdistricts with postal codes, using the official Department of Provincial Administration codes. Source: [kongvut/thai-province-data](https://github.com/kongvut/thai-province-data).

The file is generated. Do not edit it by hand.

## Capacity

All numbers below were measured, not estimated. Reproduce them with the bundled load generator:

```bash
go run ./cmd/loadtest -url http://localhost:8080 -vus 200 -n 20000 -offset 1
```

It generates unique national IDs that pass the mod-11 checksum, so every request reaches the write path instead of bouncing off format validation. Pass a comma-separated list to `-url` to spread load across replicas.

### Headline number

On the recommended spec below, with a 200 ms identity-provider call, the system sustains **~1,900 applications/second**. That is 1,000,000 applications in about 9 minutes of continuous processing, or 100,000 in under a minute.

Spread over a normal registration window the load is trivial: 1,000,000 applicants across 30 days averages 0.4 req/s. Even a first-day spike at 30× the average is ~12 req/s, roughly 0.6% of capacity.

### Connection pool sizing

12-core dev machine, 2,000 requests, unique IDs, identity stub answering instantly:

| `DB_MAX_CONNS` | 10 VU | 25 VU | 50 VU | 100 VU | 200 VU |
| --- | --- | --- | --- | --- | --- |
| 10 | 1,963 | 2,113 | 2,242 | 2,210 | 2,321 rps |
| 30 | 2,167 | 2,763 | 2,891 | 2,883 | 2,708 rps |
| 60 | 2,000 | 2,573 | 2,919 | 3,190 | 2,681 rps |

The old hardcoded value of 10 left roughly 30% of throughput on the table. Past 30 there is almost nothing left to gain.

### Sustained throughput

20,000 requests per run, `DB_MAX_CONNS=30`, identity stub instant:

| Concurrency | Throughput | p50 | p95 | p99 | Result |
| --- | --- | --- | --- | --- | --- |
| 50 | 3,026 rps | 16 ms | 21 ms | 25 ms | all 201 |
| 100 | 3,066 rps | 32 ms | 39 ms | 43 ms | all 201 |
| 200 | 3,024 rps | 65 ms | 76 ms | 85 ms | all 201 |
| 400 | 2,776 rps | 140 ms | 178 ms | 200 ms | all 201 |
| 800 | 3,044 rps | 259 ms | 286 ms | 316 ms | all 201 |

Throughput is flat from 50 to 800 concurrent while latency grows linearly — the system is saturated but degrades by queueing, never by dropping or erroring.

Data integrity after 20,000 concurrent writes: 20,000 applications, 20,000 citizens, 20,000 addresses, 20,000 verifications, zero duplicate application numbers, zero duplicate ID hashes.

### Finding the real bottleneck

Postgres in a CPU-capped container, API pinned with `GOMAXPROCS`, 200 concurrent, identity stub instant:

| Database CPU | Throughput | p50 | p95 |
| --- | --- | --- | --- |
| 1 vCPU | 1,120 rps | 189 ms | 200 ms |
| 2 vCPU | 2,323 rps | 87 ms | 95 ms |
| 4 vCPU | 2,882 rps | 68 ms | 77 ms |
| 8 vCPU | 2,888 rps | 68 ms | 77 ms |

**The database is the bottleneck, and 4 vCPU is the knee of the curve.** Going to 8 vCPU buys 0.2%.

Giving the API more CPU changes nothing: 2 vCPU produced 2,282 rps and 4 vCPU produced 2,191 rps with the database held constant — within noise of each other.

### With a realistic identity-provider call

The identity stub answers instantly, which is not realistic. With a simulated 200 ms provider call (`STUB_KYC_DELAY_MS=200`) on API 2 vCPU + database 4 vCPU:

| Concurrency | Throughput | p50 | p95 |
| --- | --- | --- | --- |
| 100 | 424 rps | 234 ms | 240 ms |
| 200 | 762 rps | 262 ms | 271 ms |
| 400 | 1,297 rps | 307 ms | 324 ms |
| 800 | 1,873 rps | 427 ms | 456 ms |

Latency here is provider wait, not database work. Throughput rises with concurrency because more calls wait in parallel.

### Replica scaling

Three API instances sharing one database, 1,200 concurrent, 200 ms provider call:

| Replicas | Pool per replica | Total connections | Throughput | p95 |
| --- | --- | --- | --- | --- |
| 1 | 16 | 16 | 2,000 rps (with errors) | 701 ms |
| 3 | 16 | 48 | 2,764 rps | 654 ms |
| 6 | 16 | 96 | 2,469 rps | 856 ms |

Scaling from 1 to 3 replicas helps. Scaling from 3 to 6 makes things worse — the bottleneck has already moved to the single database, and extra replicas just add queueing.

### The scaling trap

Postgres defaults to `max_connections = 100`. Every replica opens its own pool, so replicas multiply:

| Configuration | Connections requested | Result |
| --- | --- | --- |
| 6 replicas × pool 24 | 144 | **4,071 of 8,000 requests returned 500** |
| 10 replicas × pool 24 | 240 | **5,868 of 8,000 requests returned 500** |
| 10 replicas × pool 8 | 80 | all 8,000 succeeded, 2,595 rps |

This is the dangerous failure mode: load rises, the autoscaler adds replicas, connections exceed the limit, and the system fails harder the more it scales. Keep `maxReplicas × DB_MAX_CONNS` comfortably below `max_connections`, or put PgBouncer in front.

### Recommended production spec

| Component | Spec | Why |
| --- | --- | --- |
| API | 2 vCPU / 2 GB, **3 replicas** | One replica saturates the database; three is for availability, not throughput |
| `DB_MAX_CONNS` | **16** per replica | 3 × 16 = 48, safely under `max_connections` |
| PostgreSQL | **4 vCPU / 8 GB** | Measured knee of the curve |

Do not provision the database below 2 vCPU. At 1 vCPU throughput halves and p99 jumps from 99 ms to 309 ms — that is a cliff, not a slope.

### What these numbers do not cover

- Everything ran on one machine over loopback: no real network, no TLS, and the load generator competed for the same CPUs.
- Runs lasted seconds, not hours. Connection leaks, index bloat, and autovacuum behaviour under a growing table are untested.
- Traffic was all-success. No mix of duplicates, validation failures, or concurrent `address-search` reads.
- Concurrency above ~1,000 could not be tested: the loopback interface runs out of ephemeral ports and the load generator fails before the server does.
- The identity provider was simulated. **If the real provider rate-limits to, say, 50 req/s, then 50 req/s is the true ceiling regardless of everything above.** That number matters more than any hardware choice here.

### A burst of 100,000 or more at the same instant

Draining a simultaneous burst at ~2,500 rps takes 40 seconds for 100,000 requests and about 6 minutes 40 seconds for 1,000,000. Typical load balancer timeouts are 30–60 seconds, so a true simultaneous burst of that size fails on timeouts — and it fails at the edge, not in the database. Each in-flight request costs roughly 26 KB of memory in the API (measured: 1,000 concurrent requests held 46 MB), so a million concurrent would need about 26 GB just to hold connections open.

No amount of extra hardware fixes this, because the single database is the serialization point. The fix is to spread arrivals: stagger registration windows by the last digit of the national ID, or put a virtual waiting room in front. Staggering 1,000,000 applicants across 10 groups and an 8-hour day yields 3.5 req/s — 0.14% of capacity.

## Tests

```bash
make test
```

The repository and service tests need PostgreSQL running (`make up`); the domain and pkg tests run standalone.

The frontend has no test suite yet. Type-check it with:

```bash
cd frontend && npx vue-tsc --noEmit
```

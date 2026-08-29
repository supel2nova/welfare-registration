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

## Tests

```bash
make test
```

The repository and service tests need PostgreSQL running (`make up`); the domain and pkg tests run standalone.

The frontend has no test suite yet. Type-check it with:

```bash
cd frontend && npx vue-tsc --noEmit
```

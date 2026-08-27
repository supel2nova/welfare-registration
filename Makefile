include .env
export
MIGRATE = migrate -path backend/db/migrations -database "$(DATABASE_URL)"

.PHONY: up down migrate rollback seed api web test tidy
up:       ; docker compose up -d
down:     ; docker compose down
migrate:  ; $(MIGRATE) up
rollback: ; $(MIGRATE) down 1
seed:     ; @for f in backend/db/seed/*.sql; do echo "seed $$f"; psql "$(DATABASE_URL)" -q -f $$f; done
api:      ; cd backend && go run ./cmd/api
web:      ; cd frontend && npm run dev
test:     ; cd backend && go test ./... -count=1
tidy:     ; cd backend && go mod tidy

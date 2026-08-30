include .env
export
MIGRATE = migrate -path backend/db/migrations -database "$(DATABASE_URL)"

.PHONY: up down migrate rollback seed psql api web test tidy
up:       ; docker compose up -d --wait
down:     ; docker compose down
migrate:  ; $(MIGRATE) up
rollback: ; $(MIGRATE) down 1

seed:     ; @for f in backend/db/seed/*.sql; do echo "seed $$f"; docker compose exec -T db psql -q -v ON_ERROR_STOP=1 -U welfare -d welfare < $$f; done
psql:     ; @docker compose exec db psql -U welfare -d welfare
api:      ; cd backend && go run ./cmd/api
web:      ; cd frontend && npm run dev
test:     ; cd backend && go test ./... -count=1
tidy:     ; cd backend && go mod tidy

package testutil

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultURL = "postgres://welfare:welfare@localhost:5432/welfare?sslmode=disable"
	lockKey    = 42424201
)

var tables = []string{
	"application_status_history", "liabilities", "assets", "income_sources",
	"household_members", "applications", "identity_verifications", "citizens", "addresses",
}

func Pool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		url = os.Getenv("DATABASE_URL")
	}
	if url == "" {
		url = defaultURL
	}

	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		t.Skipf("ไม่มี postgres ให้เทส (%v) — สั่ง make up ก่อน", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		t.Skipf("ไม่มี postgres ให้เทส (%v) — สั่ง make up ก่อน", err)
	}

	ctx := context.Background()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		pool.Close()
		t.Fatalf("acquire: %v", err)
	}
	if _, err := conn.Exec(ctx, `SELECT pg_advisory_lock($1)`, lockKey); err != nil {
		conn.Release()
		pool.Close()
		t.Fatalf("advisory lock: %v", err)
	}

	t.Cleanup(func() {
		_, _ = conn.Exec(ctx, `SELECT pg_advisory_unlock($1)`, lockKey)
		conn.Release()
		pool.Close()
	})
	return pool
}

func Reset(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()

	for _, table := range tables {
		if _, err := pool.Exec(ctx, "TRUNCATE "+table+" CASCADE"); err != nil {
			t.Fatalf("truncate %s: %v", table, err)
		}
	}
	if _, err := pool.Exec(ctx, "ALTER SEQUENCE app_no_seq RESTART"); err != nil {
		t.Fatalf("reset app_no_seq: %v", err)
	}
}

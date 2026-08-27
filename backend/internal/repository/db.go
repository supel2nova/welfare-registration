package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, url string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}
	// ★ ตั้งค่าชัดเจน อย่าใช้ default เงียบๆ
	// MaxConns: เผื่อ 4 ต่อ CPU core ของ dev machine — prod ต้องวัดจริง
	cfg.MaxConns = 10
	cfg.MinConns = 2
	// อายุสั้นกว่า idle timeout ของ pgbouncer/proxy ที่จะมาในอนาคต
	cfg.MaxConnLifetime = 30 * time.Minute
	cfg.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return pool, pool.Ping(ctx)
}

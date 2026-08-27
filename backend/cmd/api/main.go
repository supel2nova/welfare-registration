package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"

	"github.com/supel2nova/welfare-registration/backend/internal/config"
	"github.com/supel2nova/welfare-registration/backend/internal/handler"
	"github.com/supel2nova/welfare-registration/backend/internal/repository"
	"github.com/supel2nova/welfare-registration/backend/internal/service"
	"github.com/supel2nova/welfare-registration/backend/internal/verifier"
	"github.com/supel2nova/welfare-registration/backend/pkg/idcrypto"
)

func main() {
	_ = godotenv.Load("../.env")
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx := context.Background()
	pool, err := repository.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer pool.Close()

	cipher, err := idcrypto.New(cfg.EncKey)
	if err != nil {
		log.Fatalf("enc: %v", err)
	}

	repo := repository.New(pool)
	apps := service.NewApplicationService(repo, verifier.Stub{}, cipher, cfg.HashPepper)
	ref := service.NewRefService(repo)

	r := handler.NewRouter(handler.Deps{
		Cfg:  cfg,
		Pool: pool,
		Repo: repo,
		Apps: apps,
		Ref:  ref,
	})

	log.Printf("listening on :%s", cfg.HTTPPort)
	if err := r.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatalf("listen: %v", err)
	}
}

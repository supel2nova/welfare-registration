package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

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

	srv := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	go func() {
		log.Printf("listening on :%s", cfg.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %v", err)
		}
	}()

	<-ctx.Done()
	stop()
	log.Println("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}

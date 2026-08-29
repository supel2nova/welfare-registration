package config

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Config struct {
	Env               string
	HTTPPort          string
	DatabaseURL       string
	CORSOrigin        string
	StubAuthEnabled   bool
	StubDefaultUserID string
	HashPepper        string
	EncKey            string
	MaxInFlight       int
	MaxBodyBytes      int64
}

func Load() Config {
	return Config{
		Env:               get("APP_ENV", "development"),
		HTTPPort:          get("HTTP_PORT", "8080"),
		DatabaseURL:       get("DATABASE_URL", ""),
		CORSOrigin:        get("CORS_ORIGIN", "http://localhost:5173"),
		StubAuthEnabled:   getBool("STUB_AUTH_ENABLED", false),
		StubDefaultUserID: get("STUB_DEFAULT_USER_ID", ""),
		HashPepper:        get("HASH_PEPPER", ""),
		EncKey:            get("ENC_KEY", ""),
		MaxInFlight:       getInt("MAX_INFLIGHT", 512),
		MaxBodyBytes:      int64(getInt("MAX_BODY_BYTES", 1<<20)),
	}
}

func (c Config) Validate() error {
	var missing []string
	for k, v := range map[string]string{
		"DATABASE_URL": c.DatabaseURL,
		"HASH_PEPPER":  c.HashPepper,
		"ENC_KEY":      c.EncKey,
	} {
		if v == "" {
			missing = append(missing, k)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		return fmt.Errorf("config: ขาด %s", strings.Join(missing, ", "))
	}
	if c.IsProduction() && c.StubAuthEnabled {
		return errors.New("config: STUB_AUTH_ENABLED=true ใน production")
	}
	return nil
}

func (c Config) IsProduction() bool { return c.Env == "production" }

func get(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func getBool(k string, def bool) bool {
	v, err := strconv.ParseBool(os.Getenv(k))
	if err != nil {
		return def
	}
	return v
}

func getInt(k string, def int) int {
	v, err := strconv.Atoi(os.Getenv(k))
	if err != nil {
		return def
	}
	return v
}

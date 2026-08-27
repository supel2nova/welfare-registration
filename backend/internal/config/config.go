package config

import (
	"os"
	"strconv"
)

type Config struct {
	Env             string
	HTTPPort        string
	DatabaseURL     string
	CORSOrigin      string
	StubAuthEnabled bool
}

func Load() Config {
	return Config{
		Env:             get("APP_ENV", "development"),
		HTTPPort:        get("HTTP_PORT", "8080"),
		DatabaseURL:     get("DATABASE_URL", ""),
		CORSOrigin:      get("CORS_ORIGIN", "http://localhost:5173"),
		StubAuthEnabled: getBool("STUB_AUTH_ENABLED", false),
	}
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

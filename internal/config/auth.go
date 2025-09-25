package config

import (
	"os"
	"time"
)

type AuthCfg struct {
	Issuer            string
	JWTSecret         string        
	AccessTTL         time.Duration 
	RefreshTTL        time.Duration 
	RefreshCookieName string        
	RefreshCookiePath string        
}

func LoadAuth() AuthCfg {
	return AuthCfg{
		JWTSecret:         getenv("JWT_SECRET", "dev-secret-change-me"),
		AccessTTL:         getenvDuration("ACCESS_TTL", 10*time.Minute),
		RefreshTTL:        getenvDuration("REFRESH_TTL", 30*24*time.Hour),
		RefreshCookieName: getenv("REFRESH_COOKIE_NAME", "rt"),
		RefreshCookiePath: getenv("REFRESH_COOKIE_PATH", "/api/auth"),
		Issuer:            getenv("ISSUER", "http://localhost:8080"),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
func getenvDuration(k string, def time.Duration) time.Duration {
	if v := os.Getenv(k); v != "" {
		d, err := time.ParseDuration(v)
		if err == nil {
			return d
		}
	}
	return def
}

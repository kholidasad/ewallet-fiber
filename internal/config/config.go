package config

import "os"

type Config struct {
	Env             string
	Port            string
	DBHost          string
	DBPort          string
	DBUser          string
	DBPass          string
	DBName          string
	JWTSecret       string
	DisplayCurrency string
	Rates           map[string]string
}

func Load() Config {
	return Config{
		Env:             getenv("APP_ENV", "development"),
		Port:            getenv("PORT", "8080"),
		DBHost:          getenv("DB_HOST", "127.0.0.1"),
		DBPort:          getenv("DB_PORT", "3306"),
		DBUser:          getenv("DB_USERNAME", "user"),
		DBPass:          getenv("DB_PASSWORD", "password"),
		DBName:          getenv("DB_NAME", "ewallet"),
		JWTSecret:       getenv("JWT_SECRET", "changeme"),
		DisplayCurrency: getenv("DISPLAY_CURRENCY", "USD"),
		Rates: map[string]string{
			"USD": getenv("RATE_USD", "1"),
			"EUR": getenv("RATE_EUR", "0.9"),
			"JPY": getenv("RATE_JPY", "155"),
		},
	}
}

func getenv(k, def string) string { if v := os.Getenv(k); v != "" { return v }; return def }

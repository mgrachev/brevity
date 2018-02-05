package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

const (
	pgConnString = "postgres://%s:%s@%s:%s/%s?sslmode=%s"

	defaultAppPort        = "8080"
	defaultAppDomain      = "http://localhost:8080"
	defaultAppTokenLength = 6
)

func AppPort() string {
	port := os.Getenv("APP_PORT")
	if port == "" {
		return defaultAppPort
	}
	return port
}

func AppEnv() string {
	return os.Getenv("APP_ENV")
}

func PGConnectionString() string {
	if url := os.Getenv("PG_CONNECTION_URL"); url != "" {
		return url
	}

	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")
	user := os.Getenv("PG_USER")
	pwd := os.Getenv("PG_PASSWORD")
	db := os.Getenv("PG_DATABASE")
	ssl := os.Getenv("PG_SSL_MODE")
	return fmt.Sprintf(pgConnString, user, pwd, host, port, db, ssl)
}

func AppDomain() string {
	domain := os.Getenv("APP_DOMAIN")
	if domain == "" {
		return defaultAppDomain
	}
	return domain
}

func AppTokenLength() int {
	length := os.Getenv("APP_TOKEN_LENGTH")
	if length == "" {
		return defaultAppTokenLength
	}

	lenInt, err := strconv.Atoi(length)
	if err != nil {
		log.Fatal(err)
	}
	return lenInt
}

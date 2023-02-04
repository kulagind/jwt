package main

import (
	"jwt/internal/services"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	services.LoadEnv()
	os.Setenv("APP_MODE", "testing")
	code := m.Run()
	os.Exit(code)
}

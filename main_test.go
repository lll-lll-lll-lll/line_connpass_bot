package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestEnv(t *testing.T) {
	wantTest := "test"
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	}
	got := os.Getenv("TEST")
	if wantTest != got {
		t.Error("failure", got)
	}

}

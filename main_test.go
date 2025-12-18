package main

import (
	"os"
	"testing"
)

func TestLoadEnvFromFile_NotFound(t *testing.T) {
	if err := loadEnvFromFile("does_not_exist"); err == nil {
		t.Fatal("expected file not found error, got nil")
	}
}

func TestLoadEnvFromFile_Success(t *testing.T) {
	file, err := os.CreateTemp("", "testenv")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())
	file.Close()

	fileContent := []byte("HOST=127.0.0.1\nPORT=8080\n")
	if err := os.WriteFile(file.Name(), fileContent, 0600); err != nil {
		t.Fatalf("failed to write env content: %v", err)
	}

	os.Unsetenv("HOST")
	os.Unsetenv("PORT")

	if err := loadEnvFromFile(file.Name()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got := os.Getenv("HOST"); got != "127.0.0.1" {
		t.Fatalf("HOST loaded incorrectly: got %v", got)
	}

	if got := os.Getenv("PORT"); got != "8080" {
		t.Fatalf("PORT loaded incorrectly: got %v", got)
	}
}

func TestLoadEnvFromFile_Missing(t *testing.T) {
	file, err := os.CreateTemp("", "testenv")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())
	file.Close()

	fileContent := []byte("HOST=127.0.0.1\n")
	if err := os.WriteFile(file.Name(), fileContent, 0600); err != nil {
		t.Fatalf("failed to write env content: %v", err)
	}

	os.Unsetenv("HOST")
	os.Unsetenv("PORT")

	if err := loadEnvFromFile(file.Name()); err == nil {
		t.Fatal("expected error for missing PORT, got nil")
	}
}

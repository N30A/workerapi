package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var requiredEnvVars = []string{"HOST", "PORT"}

func loadEnvFromFile(filenames ...string) error {
	err := godotenv.Load(filenames...)
	if err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	missing := make([]string, 0, len(requiredEnvVars))
	for _, key := range requiredEnvVars {
		if os.Getenv(key) == "" {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("required env variable(s) are missing: %s", strings.Join(missing, ", "))
	}

	return nil
}

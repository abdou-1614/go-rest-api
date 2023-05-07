package common

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	prod := os.Getenv("Prod")

	if prod != "true" {
		err := godotenv.Load()

		if err != nil {
			return err
		}
	}

	return nil
}

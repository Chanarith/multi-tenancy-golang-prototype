package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct{}

func (e Env) Initialize() {
	if os.Getenv("GIN_MODE") != "release" {
		err := godotenv.Load(".env.development")
		if err != nil {
			log.Fatal("Error loading environment variables")
		}
	}
}

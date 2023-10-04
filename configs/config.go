package configs

import (
	"fmt"
	"github.com/joho/godotenv"
)

func Load() {
	// Load env variables
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("error: cannot find .env file in the project root")
	}

	//	TODO Setup Database connection.
}

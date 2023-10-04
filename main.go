package main

import (
	"demerzel-badges/api"
	"demerzel-badges/configs"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	configs.Load()

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to conver PORT to integer: %v", err))
	}

	server := api.NewServer(uint16(port), api.SetupRoutes())
	server.Listen()
}

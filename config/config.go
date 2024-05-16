package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	ServerPort string
	// Others can be added here
}

func LoadConfig() Config {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	} else {
		if _, err := strconv.Atoi(port); err != nil {
			log.Fatalf("Invalid port number: %s", port)
		}
	}

	return Config{
		ServerPort: port,
	}
}

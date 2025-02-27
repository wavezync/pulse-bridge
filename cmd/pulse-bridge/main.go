package main

import (
	"fmt"
	"net/http"
	"os"
	"wavezync/pulse-bridge/internal/config"
	"wavezync/pulse-bridge/internal/handler"
)

func main() {
	configPath := "config.yml"
	configData, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		os.Exit(1)
	}

	cfg, err := config.Init(string(configData))
	if err != nil {
		fmt.Printf("Error initializing config: %v\n", err)
		os.Exit(1)
	}

	http.HandleFunc("/health", handler.Health)
	http.HandleFunc("/status", handler.Status)

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(*cfg)
}

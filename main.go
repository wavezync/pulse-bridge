package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"wavezync/pulse-bridge/internal/config"
	"wavezync/pulse-bridge/service/health"
	"wavezync/pulse-bridge/service/status"
)

//go:embed config.yml
var configFile string

func main() {
	cfg, err := config.Init(configFile)

	if err != nil {
		fmt.Println(err)
	}

	http.HandleFunc("/health", health.GetHealth)
	http.HandleFunc("/status", status.GetStatus)

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(*cfg)
}

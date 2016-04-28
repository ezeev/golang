package kolektor

import (
	"encoding/json"
	"fmt"
	"os"
)

// Configuration Struct to hold config information
type Configuration struct {
	ListenPort  string
	Backend     string
	Debug       bool
	BackendArgs map[string]string
	Interval    float64
}

// LoadConfig Loads config information from json into the Configuration struct
func LoadConfig(path string) Configuration {
	file, _ := os.Open(path)
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("Error loading config:", err)
	}
	return configuration
}

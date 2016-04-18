package kolektor

import (
    "encoding/json"
    "os"
    "fmt"
)

type Configuration struct {
    ListenPort string
    Backend   string
    Debug     bool
    BackendArgs map[string]string
    Interval float64
}

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

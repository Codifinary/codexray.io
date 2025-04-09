package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

func SaveToJSON(serviceName string, data interface{}) error {

	logsDir := filepath.Join("logs", serviceName)
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %v", err)
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	randomNum := rand.Intn(10000)
	filename := filepath.Join(logsDir, fmt.Sprintf("%s_%04d.json", timestamp, randomNum))

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	log.Printf("Saved data to %s", filename)
	return nil
}

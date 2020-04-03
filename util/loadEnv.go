package util

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// LoadEnv takes the path of a .env file and uses it to set environment variables
func LoadEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Unable to open env file at", path, ":", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		data := strings.Split(line, "=")
		key := strings.TrimSpace(data[0])
		value := strings.TrimSpace(data[1])

		os.Setenv(key, value)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error reading env file at", path, ":", err)
	}
}

// LoadEnvFromPath will load the environment variables
func LoadEnvFromPath(dirPath string) {
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if filepath.Ext(path) == ".env" {
				LoadEnv(path)
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal("Error reading .env files at path", dirPath, ":", err)
	}
}

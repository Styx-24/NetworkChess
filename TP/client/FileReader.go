package client

import (
	"io"
	"log"
	"os"
)

func GetServerAddress() string {

	file, err := os.Open(".\\client\\config.txt")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	return string(content)
}

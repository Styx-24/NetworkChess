package server

import (
	"io"
	"log"
	"os"
	"strings"
)

func GetConfig() (string, string) {

	file, err := os.Open(".\\server\\config.txt")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	config := strings.Split(string(content), "\n")

	return strings.Replace(config[0], "\r", "", -1), config[1]
}

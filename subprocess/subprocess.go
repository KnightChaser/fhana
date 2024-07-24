// <PROJECT_DIR>/subprocess/subprocess.go
package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

// Generates a random hexadecimal string for encryption key
func generateRandomHex(bits int) (string, error) {
	if bits%8 != 0 {
		log.Fatalf("Invalid number of bits: %d, must be divisible by 8", bits)
		return "", fmt.Errorf("Invalid number of bits: %d, must be divisible by 8", bits)
	}
	bytes := make([]byte, bits/8)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Saves the given stringified data to a specified file
func saveStringToFile(filename string, data string) error {
	return os.WriteFile(filename, []byte(data), 0644)
}

func main() {
	randomHexString, err := generateRandomHex(128)
	if err != nil {
		log.Fatalf("Error generating random hex: %v", err)
	}

	fmt.Printf("Random hex: %s\n", randomHexString)

	filename := "key.txt"
	err = saveStringToFile(filename, randomHexString)
	if err != nil {
		log.Fatalf("Error saving to file: %v", err)
	}

	fmt.Printf("Random hex saved to %s\n", filename)
}

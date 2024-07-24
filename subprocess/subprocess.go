// <PROJECT DIR>/subprocess/subprocess.go
package main

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"

	"fhana/subprocess/ciphers"

	"github.com/joho/godotenv"
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
	// Parse the --decrypt flag
	decrypt := flag.Bool("decrypt", false, "Decrypt files instead of encrypting")
	flag.Parse()

	// Load the environment variables from .env file to specify the target directory
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
		return
	}
	targetDirectoryAbsolutePath := os.Getenv("TARGET_ABSOLUTE_DIRPATH")
	if targetDirectoryAbsolutePath == "" {
		log.Fatalf("TARGET_ABSOLUTE_DIRPATH not found in .env file")
		return
	}

	if *decrypt {
		// Decryption procedure is requested. Load the AES key from the file and decrypt the directory
		// Load the AES key from the file
		keyHex, err := os.ReadFile("key.txt")
		if err != nil {
			log.Fatalf("Error reading key file: %v", err)
			return
		}

		aesDecryptionKey, err := hex.DecodeString(string(keyHex))
		if err != nil {
			log.Fatalf("Error decoding key: %v", err)
			return
		}

		err = ciphers.DecryptDirectory(targetDirectoryAbsolutePath, aesDecryptionKey)
		if err != nil {
			log.Fatalf("Error decrypting directory: %v", err)
			return
		}

		fmt.Println("Decryption completed successfully")
	} else {
		// Encryption procedure is requested. Generate a random AES key and encrypt the directory
		// Generate AES key file
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

		// Start encryption to the target directory
		aesEncryptionKey, err := hex.DecodeString(randomHexString)
		if err != nil {
			log.Fatalf("Error decoding random hex: %v", err)
			return
		}

		err = ciphers.EncryptDirectory(targetDirectoryAbsolutePath, aesEncryptionKey)
		if err != nil {
			log.Fatalf("Error encrypting directory: %v", err)
			return
		}

		fmt.Println("Encryption completed successfully")
	}
}

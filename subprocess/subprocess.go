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
	// Parse flags
	encrypt := flag.Bool("encrypt", false, "Encrypt files in the target directory")
	decrypt := flag.Bool("decrypt", false, "Decrypt files instead of encrypting")
	targetDir := flag.String("target-directory", "", "The target directory to encrypt/decrypt")
	keyFile := flag.String("key-file", "key.txt", "The file to save/load the encryption key")
	flag.Parse()

	if *targetDir == "" {
		log.Fatalf("Error: target directory must be specified")
		return
	}

	if *decrypt {
		// Decryption procedure is requested. Load the AES key from the file and decrypt the directory
		keyHex, err := os.ReadFile(*keyFile)
		if err != nil {
			log.Fatalf("Error reading key file: %v", err)
			return
		}

		aesDecryptionKey, err := hex.DecodeString(string(keyHex))
		if err != nil {
			log.Fatalf("Error decoding key: %v", err)
			return
		}

		err = ciphers.DecryptDirectory(*targetDir, aesDecryptionKey)
		if err != nil {
			log.Fatalf("Error decrypting directory: %v", err)
			return
		}

		fmt.Println("Decryption completed successfully")
	} else if *encrypt {
		// Encryption procedure is requested. Generate a random AES key and encrypt the directory
		randomHexString, err := generateRandomHex(128)
		if err != nil {
			log.Fatalf("Error generating random hex: %v", err)
		}

		fmt.Printf("Random hex: %s\n", randomHexString)

		err = saveStringToFile(*keyFile, randomHexString)
		if err != nil {
			log.Fatalf("Error saving to file: %v", err)
		}

		fmt.Printf("Random hex saved to %s\n", *keyFile)

		aesEncryptionKey, err := hex.DecodeString(randomHexString)
		if err != nil {
			log.Fatalf("Error decoding random hex: %v", err)
			return
		}

		err = ciphers.EncryptDirectory(*targetDir, aesEncryptionKey)
		if err != nil {
			log.Fatalf("Error encrypting directory: %v", err)
			return
		}

		fmt.Println("Encryption completed successfully")
	} else {
		log.Fatalf("Error: either --encrypt or --decrypt flag must be specified")
		return
	}
}

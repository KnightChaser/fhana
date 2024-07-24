// <PROJECT DIR>/subprocess/ciphers/decryption.go
// Use this code in case you need to decrypt files encrypted using AES encryption algorithm
package ciphers

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Decrypt the given file (specified by filepath) in AES using the given key
func decryptFileAES(filepath string, key []byte) error {
	ciphertext, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Error reading encrypted file: %v", err)
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

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Fatalf("Error decrypting file: %v", err)
		return err
	}

	decryptedFilepath := strings.TrimSuffix(filepath, ".knightz.enc")
	err = os.WriteFile(decryptedFilepath, plaintext, 0644)
	if err != nil {
		log.Fatalf("Error writing decrypted file: %v", err)
		return err
	}

	return nil
}

// DecryptDirectory walks through the directory and decrypts each encrypted file
func DecryptDirectory(targetDirectory string, key []byte) error {
	err := filepath.Walk(targetDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".knightz.enc") {
			fmt.Printf("Decrypting file: %s\n", path)
			err = decryptFileAES(path, key)
			if err != nil {
				log.Fatalf("Error decrypting file: %v", err)
				return err
			}
			fmt.Printf("Decrypted file: %s\n", path)
		}
		return nil
	})

	return err
}

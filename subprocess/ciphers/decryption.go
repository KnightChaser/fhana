// <PROJECT DIR>/subprocess/ciphers/decryption.go
package ciphers

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
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

	decryptedFilepath := strings.TrimSuffix(filepath, ".knightz.encrypted")
	err = os.WriteFile(decryptedFilepath, plaintext, 0644)
	if err != nil {
		log.Fatalf("Error writing decrypted file: %v", err)
		return err
	}

	// Delete the encrypted file
	err = os.Remove(filepath)
	if err != nil {
		log.Fatalf("Error deleting encrypted file: %v", err)
		return err
	}

	return nil
}

// A worker reads the file (given by its filepath) from the channel, then conduct decryption
// This will be executed concurrently by goroutines
func decryptionWorker(workerNumber int, filepaths <-chan string, key []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	for filepath := range filepaths {
		err := decryptFileAES(filepath, key)
		if err != nil {
			fmt.Printf("[Worker #%02d]Decrypting file: %s => failed (reason: %v)\n", workerNumber, filepath, err)
			continue
		} else {
			fmt.Printf("[Worker #%02d]Decrypting file: %s => success\n", workerNumber, filepath)
		}
	}
}

// DecryptDirectory walks through the directory and decrypts each encrypted file using multiple goroutines
func DecryptDirectory(targetDirectory string, key []byte) error {
	// Create a channel to share filepaths with the decryption worker goroutines
	targetFiles := make(chan string, 1024)
	var wg sync.WaitGroup

	const numberOfDecryptionWorkers = 10
	for i := 0; i < numberOfDecryptionWorkers; i++ {
		wg.Add(1)
		go decryptionWorker(i+1, targetFiles, key, &wg)
	}

	// Walk through the target directory and send the filepaths to the channel
	err := filepath.Walk(targetDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Encrypted file by this program("fhana") has ".knightz.encrypted" suffix to the original filename
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".knightz.encrypted") {
			targetFiles <- path
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking through directory: %v", err)
		return err
	}

	// All target files have been sent to the channel
	// Close the files channel and wait for all workers to finish
	close(targetFiles)

	// Wait for all workers to finish
	wg.Wait()

	return nil
}

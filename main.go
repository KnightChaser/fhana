// <PROJECT DIR>/main.go
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os/exec"
)

func main() {
	// Parse the flags
	encrypt := flag.Bool("encrypt", false, "Encrypt files in the target directory")
	decrypt := flag.Bool("decrypt", false, "Decrypt files instead of encrypting")
	targetDir := flag.String("target-directory", "", "The target directory to encrypt/decrypt")
	keyFile := flag.String("key-file", "key.txt", "The file to save/load the encryption key")
	flag.Parse()

	if *targetDir == "" {
		fmt.Println("Error: target directory must be specified")
		return
	}

	var cmd *exec.Cmd
	if *decrypt {
		cmd = exec.Command("go", "run", "./subprocess/subprocess.go", "--decrypt", "--target-directory", *targetDir, "--key-file", *keyFile)
	} else if *encrypt {
		cmd = exec.Command("go", "run", "./subprocess/subprocess.go", "--encrypt", "--target-directory", *targetDir, "--key-file", *keyFile)
	} else {
		fmt.Println("Error: either --encrypt or --decrypt flag must be specified")
		return
	}

	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: %v: %v\n", err, stderr.String())
		return
	}

	fmt.Println("Subprocess output:")
	fmt.Println(out.String())

	fmt.Println("Subprocess ran successfully")
}

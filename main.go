// <PROJECT DIR>/main.go
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os/exec"
)

func target() {
	fmt.Println("This is the target function")
}

func main() {
	// Parse the --decrypt flag
	decrypt := flag.Bool("decrypt", false, "Decrypt files instead of encrypting")
	flag.Parse()

	target()

	var cmd *exec.Cmd
	if *decrypt {
		cmd = exec.Command("go", "run", "./subprocess/subprocess.go", "--decrypt")
	} else {
		cmd = exec.Command("go", "run", "./subprocess/subprocess.go")
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

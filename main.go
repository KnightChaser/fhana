// <PROJECT_DIR>/main.go
package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

func target() {
	fmt.Println("This is the target function")
}

func main() {
	target()

	cmd := exec.Command("go", "run", "./subprocess/subprocess.go")

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: %v: %v\n", err, stderr.String())
		return
	}

	fmt.Printf("Subprocess output: %v\n", out.String())
	fmt.Println("Subprocess ran successfully")
}

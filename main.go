// <PROJECT_DIR>/main.go
package main

import (
	"fmt"
	"os/exec"
)

func target() {
	fmt.Println("This is the target function")
}

func main() {
	target()

	cmd := exec.Command("go", "run", "./subprocess/subprocess.go")

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Subprocess ran successfully")
}

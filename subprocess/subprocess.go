// <PROJECT_DIR>/subprocess/subprocess.go
package main

import (
	"fmt"
)

func target() {
	fmt.Println("This is the target function in the subprocess")
}

func main() {
	target()
}

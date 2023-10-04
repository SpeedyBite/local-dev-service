package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// Define the command and its arguments
	command := "docker"
	args := []string{"run", "-it", "--rm", "f7f18a3e6906", "/bin/bash"}

	// Create a new command and attach it to the standard input/output
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		fmt.Println("Error running command:", err)
	}
}

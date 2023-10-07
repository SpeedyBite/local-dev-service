package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// Create a channel to receive key presses
	keyCh := make(chan rune)

	// Start a goroutine to read key presses from stdin
	go func() {
		fmt.Println("Press 'q' to quit")
		reader := bufio.NewReader(os.Stdin)
		for {
			// Check if there is data available to read
			if reader.Buffered() > 0 {
				key, _, err := reader.ReadRune()
				if err != nil {
					close(keyCh)
					return
				}
				if key == 'q' {
					close(keyCh)
					return
				}
				keyCh <- key
			}
		}
	}()

	// Print each key press received from the channel
	for key := range keyCh {
		fmt.Printf("Key pressed: %c\n", key)
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Create a scanner to read an input
	scanner := bufio.NewScanner(os.Stdin)

	// Start an infinite loop for the REPL
	for {
		// Print the prompt without a newline
		fmt.Print("Pokedex > ")

		// Wait for user input
		scanner.Scan()

		// Get the text that the user entered
		input := scanner.Text()

		// Clean the input
		cleaned := cleanInput(input)
		if len(cleaned) == 0 {
			continue
		}

		// Get the 1st word as the command
		command := cleaned[0]
		fmt.Printf("Your command was: %s\n", command)
	}
}

func cleanInput(text string) []string {
	// Convert to lowercase
	text = strings.ToLower(text)

	// Split by whitespace and filter out empty strings
	words := strings.Fields(text)

	return words
}

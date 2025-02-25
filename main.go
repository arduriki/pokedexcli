package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

// Commands as a global variable, but not initialized
var commands map[string]cliCommand

func main() {
	// Initialize commands
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
	}
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
		commandName := cleaned[0]

		// Look up the command in the registry
		command, exists := commands[commandName]

		if exists {
			// If command exists, call its callback
			err := command.callback()
			if err != nil {
				fmt.Printf("Error executing command: %s\n", err)
			}
		} else {
			fmt.Println("Unknown command")
		}

	}
}

func cleanInput(text string) []string {
	// Convert to lowercase
	text = strings.ToLower(text)

	// Split by whitespace and filter out empty strings
	words := strings.Fields(text)

	return words
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	// Get all command names
	names := make([]string, 0, len(commands))
	for name := range commands {
		names = append(names, name)
	}

	// Sort command names alphabetically
	sort.Strings(names)

	// Print commands in alphabetical order
	for _, name := range names {
		cmd := commands[name]
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}

	return nil
}

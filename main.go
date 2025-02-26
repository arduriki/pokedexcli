package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/arduriki/pokedexcli/internal/pokecache"
)

type Config struct {
	Next          string
	Previous      string
	Cache         *pokecache.Cache
	CaughtPokemon map[string]Pokemon
}

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, []string) error
}

// LocationAreasResp models the response from the location areas API
type LocationAreasResp struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// LocationAreaDetails models the response from the location area details API
type LocationAreaDetails struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance   int `json:"chance"`
				MaxLevel int `json:"max_level"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

// Pokemon represents a caught Pokemon
type Pokemon struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

var commands map[string]cliCommand

func main() {
	// Initialize the Config with base URL
	config := &Config{
		Next:          "https://pokeapi.co/api/v2/location-area/",
		Previous:      "",
		Cache:         pokecache.NewCache(5 * time.Minute),
		CaughtPokemon: make(map[string]Pokemon),
	}

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
		"map": {
			name:        "map",
			description: "Display the next 20 location areas in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 location areas in the Pokemon world",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Explore a location area to find Pokemon. Usage: explore [location-area-name]",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a Pokemon. Usage: catch [pokemon-name]",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a caught Pokemon. Usage: inspect [pokemon-name]",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List all the Pokemon in your Pokedex",
			callback:    commandPokedex,
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
			// Get the arguments
			args := []string{}
			if len(cleaned) > 1 {
				args = cleaned[1:]
			}

			// If command exists, call its callback with config
			err := command.callback(config, args)
			if err != nil {
				fmt.Println(err)
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

func commandExit(cfg *Config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *Config, args []string) error {
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

func commandMap(cfg *Config, args []string) error {
	// Check if we have a URL to fetch
	if cfg.Next == "" {
		return fmt.Errorf("no more locations to fetch")
	}

	// Get location areas from the API
	locationsResp, err := getLocationAreas(cfg.Next, cfg.Cache)
	if err != nil {
		return err
	}

	// Update the Next and Previous URLs in the config
	if locationsResp.Next != nil {
		cfg.Next = *locationsResp.Next
	} else {
		cfg.Next = ""
	}

	if locationsResp.Previous != nil {
		cfg.Previous = *locationsResp.Previous
	} else {
		cfg.Previous = ""
	}

	// Display the location areas
	for _, location := range locationsResp.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapb(cfg *Config, args []string) error {
	// Check if we have a previous URL to fetch
	if cfg.Previous == "" {
		fmt.Println("You're on the first page.")
		return nil
	}

	// Get location areas from the API
	locationsResp, err := getLocationAreas(cfg.Previous, cfg.Cache)
	if err != nil {
		return err
	}

	// Update the next and previous URLs in the config
	if locationsResp.Next != nil {
		cfg.Next = *locationsResp.Next
	} else {
		cfg.Next = ""
	}

	if locationsResp.Previous != nil {
		cfg.Previous = *locationsResp.Previous
	} else {
		cfg.Previous = ""
	}

	// Display the location areas
	for _, location := range locationsResp.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandExplore(cfg *Config, args []string) error {
	// Check if the location area name was provided
	if len(args) == 0 {
		return fmt.Errorf("please provide a location area name to explore")
	}

	// Get the location area name from the arguments
	locationName := args[0]

	// Construct the URL for the area details
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", locationName)

	// Get location area details from the API
	locationDetails, err := getLocationAreaDetails(url, cfg.Cache)
	if err != nil {
		return err
	}

	// Display the location area name
	fmt.Printf("exploring %s...\n", locationDetails.Name)

	// Check if there are any Pokemon encounters
	if len(locationDetails.PokemonEncounters) == 0 {
		fmt.Println("No Pokemon found in this area.")
		return nil
	}

	// Display the Pokemon found in this area
	fmt.Println("Found Pokemon:")
	for _, encounter := range locationDetails.PokemonEncounters {
		fmt.Printf("- %s\n", encounter.Pokemon.Name)
	}

	return nil
}

func getLocationAreaDetails(url string, cache *pokecache.Cache) (LocationAreaDetails, error) {
	// Make the HTTP GET request
	body, err := makeGetRequest(url, cache)
	if err != nil {
		return LocationAreaDetails{}, err
	}

	// Unmarshal the JSON into the struct
	var locationDetails LocationAreaDetails
	err = json.Unmarshal(body, &locationDetails)
	if err != nil {
		return LocationAreaDetails{}, err
	}

	return locationDetails, nil
}

func makeGetRequest(url string, cache *pokecache.Cache) ([]byte, error) {
	// Try to get from cache first
	if cachedData, ok := cache.Get(url); ok {
		// If found in cache, return it immediately
		return cachedData, nil
	}

	// Create an HTTP client
	client := &http.Client{}

	// Create a GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// Close the response body once it's done
	defer resp.Body.Close()

	// Check if the response status code is OK (200)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Save to cache before returning
	cache.Add(url, body)

	return body, nil
}

func getLocationAreas(url string, cache *pokecache.Cache) (LocationAreasResp, error) {
	// Make the HTTP GET request
	body, err := makeGetRequest(url, cache)
	if err != nil {
		return LocationAreasResp{}, err
	}

	// Unmarshal the JSON into our struct
	var locationsResp LocationAreasResp
	err = json.Unmarshal(body, &locationsResp)
	if err != nil {
		return LocationAreasResp{}, err
	}

	return locationsResp, nil
}

func commandCatch(cfg *Config, args []string) error {
	// Check if the Pokemon name was provided
	if len(args) == 0 {
		return fmt.Errorf("Please provide a Pokemon name to catch")
	}

	// Get the Pokemon name from the arguments
	pokemonName := args[0]

	// Print the throwing message
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	// Construct the URL for the Pokemon API
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", pokemonName)

	// Get Pokemon details from the API
	pokemonDetails, err := getPokemonDetails(url, cfg.Cache)
	if err != nil {
		return fmt.Errorf("failed to find Pokemon: %s", err)
	}

	// Calculate catch probability (lower base_experience = higher chance)
	catchProbability := calculateCatchProbability(pokemonDetails.BaseExperience)

	// Generate a random number
	randomValue := rand.Float64()

	// Check if Pokemon is caught
	if randomValue <= catchProbability {
		fmt.Printf("%s was caught!\n", pokemonDetails.Name)

		// Add to caught Pokemon map
		if cfg.CaughtPokemon == nil {
			cfg.CaughtPokemon = make(map[string]Pokemon)
		}
		cfg.CaughtPokemon[pokemonDetails.Name] = pokemonDetails

		fmt.Println("You may now inspect it with the inspect command")
	} else {
		fmt.Printf("%s escaped!\n", pokemonDetails.Name)
	}

	return nil
}

func calculateCatchProbability(baseExperience int) float64 {
	// Simple formula that gives lower chance for higher baseExperience
	baseProb := 0.7
	expFactor := float64(baseExperience) / 200.0

	// Clamp probability between 0.1 and 0.9
	probability := math.Max(0.1, math.Min(0.9, baseProb-expFactor))

	return probability
}

func getPokemonDetails(url string, cache *pokecache.Cache) (Pokemon, error) {
	// Make the HTTP GET request
	body, err := makeGetRequest(url, cache)
	if err != nil {
		return Pokemon{}, err
	}

	// Unmarshal the JSON into the struct
	var pokemonDetails Pokemon
	err = json.Unmarshal(body, &pokemonDetails)
	if err != nil {
		return Pokemon{}, err
	}

	return pokemonDetails, nil
}

func commandInspect(cfg *Config, args []string) error {
	// Check if the Pokemon name was provided
	if len(args) == 0 {
		return fmt.Errorf("please provide a Pokemon name to inspect")
	}

	// Get the Pokemon name from the arguments
	pokemonName := args[0]

	// Check if the Pokemon is in the caught map
	pokemon, ok := cfg.CaughtPokemon[pokemonName]
	if !ok {
		return fmt.Errorf("you haven't caught that pokemon")
	}

	// Display the Pokemon details
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)

	// Display stats
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf(" - %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	// Display types
	fmt.Println("Types:")
	for _, typeInfo := range pokemon.Types {
		fmt.Printf(" - %s\n", typeInfo.Type.Name)
	}

	return nil
}

func commandPokedex(cfg *Config, args []string) error {
	// Check if the Pokedex is empty
	if len(cfg.CaughtPokemon) == 0 {
		fmt.Println("Your Pokedex is empty. Go catch some Pokemon!")
		return nil
	}

	// Display the Pokedex
	fmt.Println("Your Pokedex:")
	for name := range cfg.CaughtPokemon {
		fmt.Printf(" - %s\n", name)
	}

	return nil
}

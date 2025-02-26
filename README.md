# PokedexCLI

A Pokémon encyclopedia in your terminal! This interactive command-line application allows you to explore the Pokémon world, discover Pokémon in different location areas, catch them, and build your own Pokédex.

![Pokedex](https://img.shields.io/badge/Pokedex-CLI-red)
![Go](https://img.shields.io/badge/Go-1.24+-00ADD8)

## Features

- **Explore the Pokémon world**: Navigate through different location areas
- **Discover Pokémon**: Find which Pokémon inhabit each area
- **Catch 'em all**: Try your luck at catching Pokémon (higher-level Pokémon are harder to catch!)
- **Build your Pokédex**: Keep track of all the Pokémon you've caught
- **Inspect Pokémon**: View detailed information about your caught Pokémon
- **Caching system**: Efficient data management with timed cache to reduce API calls

## Installation

### Prerequisites

- Go 1.24 or higher

### Building from source

1. Clone the repository:
   ```
   git clone https://github.com/arduriki/pokedexcli.git
   ```

2. Navigate to the project directory:
   ```
   cd pokedexcli
   ```

3. Build the application:
   ```
   go build
   ```

4. Run the application:
   ```
   ./pokedexcli
   ```

## How to Play

1. **Start the application** by running `./pokedexcli` in your terminal
2. You'll be greeted with a `Pokedex >` prompt
3. Type `help` to see all available commands
4. Use `map` to view location areas you can explore
5. Use `explore [location-name]` to see which Pokémon are in that area
6. Try to catch Pokémon with `catch [pokemon-name]`
7. View your caught Pokémon with `pokedex`
8. Inspect a caught Pokémon's details with `inspect [pokemon-name]`

## Commands

| Command | Description | Example |
|---------|-------------|---------|
| `help` | Displays all available commands | `help` |
| `exit` | Exits the Pokedex application | `exit` |
| `map` | Shows the next 20 location areas | `map` |
| `mapb` | Shows the previous 20 location areas | `mapb` |
| `explore` | Lists Pokémon in a specified area | `explore pallet-town` |
| `catch` | Attempts to catch a Pokémon | `catch pikachu` |
| `inspect` | Shows details of a caught Pokémon | `inspect pikachu` |
| `pokedex` | Lists all caught Pokémon | `pokedex` |

## Example Session

```
Pokedex > map
canalave-city-area
eterna-city-area
pastoria-city-area
sunyshore-city-area
sinnoh-pokemon-league-area
...

Pokedex > explore eterna-city-area
exploring eterna-city-area...
Found Pokemon:
- glameow
- stunky
- bronzor
- rapidash

Pokedex > catch glameow
Throwing a Pokeball at glameow...
glameow was caught!
You may now inspect it with the inspect command

Pokedex > inspect glameow
Name: glameow
Height: 5
Weight: 39
Stats:
 - hp: 49
 - attack: 55
 - defense: 42
 - special-attack: 42
 - special-defense: 37
 - speed: 85
Types:
 - normal

Pokedex > pokedex
Your Pokedex:
 - glameow
```

## How It Works

This application uses the [PokeAPI](https://pokeapi.co/) to fetch data about Pokémon and location areas. It implements a simple caching system to reduce API calls and improve performance.

When catching Pokémon, the success rate is influenced by the Pokémon's base experience - powerful Pokémon are harder to catch!

## Future Enhancements

- Update the CLI to support the "up" arrow to cycle through previous commands
- Simulate battles between Pokémon
- Add more unit tests
- Refactor code for better organization and testability
- Keep Pokémon in a "party" and allow them to level up
- Allow caught Pokémon to evolve after a set amount of time
- Persist a user's Pokédex to disk to save progress between sessions
- Make exploration more interactive with directional choices
- Implement random encounters with wild Pokémon
- Add support for different types of Poké Balls with varying catch rates

## Contributing

Contributions are welcome! Feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Data provided by [PokeAPI](https://pokeapi.co/)
- Inspired by the Pokémon game series
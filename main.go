package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Adfmu1/pokedex/internal/pokecache"
)

var commands map[string]cliCommand

func main() {
	fmt.Println("Welcome to the Pokedex!")
	initCommands()
	scanner := bufio.NewScanner(os.Stdin)

	var conf config
	conf.Cache = pokecache.NewCache(5)
	conf.NextUrl = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"

	commandHelp(&conf)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := cleanInput(scanner.Text())
		if len(input) == 0 {
			continue
		}
		command := input[0]

		if i, ok := commands[command]; !(ok) {
			fmt.Println("Unknown command")
		} else {
			if i.name == "explore" {
				err := i.callback(&conf, input[1])
				if err != nil {
					fmt.Printf("Error occured while trying to call %s command: %v", command, err)
				}
			} else {
				err := i.callback(&conf)
				if err != nil {
					fmt.Printf("Error occured while trying to call %s command: %v", command, err)
				}
			}

		}
	}
}

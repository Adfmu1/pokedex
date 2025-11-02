package main

import (
	"bufio"
	"os"
	"fmt"
	)

var commands map[string]cliCommand

func main() {
	fmt.Println("Welcome to the Pokedex!")
	initCommands()
	scanner := bufio.NewScanner(os.Stdin)

	var conf config
	conf.NextUrl= "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"


	commandHelp(&conf)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := cleanInput(scanner.Text())
		if len(input) == 0 { continue }
		command := input[0]

		if i, ok := commands[command]; !(ok) {
			fmt.Println("Unknown command")
		} else {
			err := i.callback(&conf)
			if err != nil {
				fmt.Println(fmt.Sprintf("Error occured while trying to call %s command: %v", command, err))
			}
		}
	}
}
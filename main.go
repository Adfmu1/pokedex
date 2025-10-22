package main

import (
	"bufio"
	"os"
	"fmt"
	"strings"
	)

type cliCommand struct {
	name			string
	description		string
	callback		func() error
}

var commands map[string]cliCommand

func initCommands() {
	commands = map[string]cliCommand {
		"exit": {
			name: 			"exit",
			description:	"Exit the Pokedex",
			callback:		commandExit,
		},
		"help": {
			name:			"help",
			description:	"Displays a help message",
			callback:		commandHelp,
		},
	}
}

func cleanInput(text string) []string {
	str := strings.Fields(text)
	for i := 0; i < len(str); i++ {
		str[i] = strings.ToLower(str[i])
	}

	return str
}

func commandExit() error {
	defer os.Exit(0)
	fmt.Println("Closing the Pokedex... Goodbye!")
	return nil
}

func commandHelp() error {
	fmt.Println("Usage:\n")
	for _, comm := range commands {
		fmt.Println(comm.name, ":", comm.description)
	}
	return nil
}

func main() {
	fmt.Println("Welcome to the Pokedex!")
	initCommands()
	scanner := bufio.NewScanner(os.Stdin)

	commandHelp()

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := cleanInput(scanner.Text())
		command := input[0]

		if i, ok := commands[command]; !(ok) {
			fmt.Println("Unknown command")
		} else {
			err := i.callback()
			if err != nil {
				fmt.Println(fmt.Sprintf("Error occured while trying to call %s command: %v", command, err))
			}
		}
	}
}
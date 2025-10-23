package main

import (
	"bufio"
	"os"
	"fmt"
	"strings"
	"net/http"
	"io"
	"encoding/json"
	)

type cliCommand struct {
	name			string
	description		string
	callback		func(conf *config) error
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
		"map": {
			name:			"map",
			description:	"Displays the 20 names of location areas, call again to display more",
			callback:		commandMap,
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

func commandExit(conf *config) error {
	defer os.Exit(0)
	fmt.Println("Closing the Pokedex... Goodbye!")
	return nil
}

func commandHelp(conf *config) error {
	fmt.Println("Usage:\n")
	for _, comm := range commands {
		fmt.Println(comm.name, ":", comm.description)
	}
	return nil
}

func commandMap(conf *config) error {
	url := conf.NextUrl
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Encountered an error while trying to get url:", err)
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if res.StatusCode > 299 {
		fmt.Println("Bad status code:", res.StatusCode)
		return err
	}
	if err != nil {
		fmt.Println("Encountered an error while retrieving body:", err)
		return err
	}

	mLocation := LocationAreas{}
	err = json.Unmarshal(body, &mLocation)
	if err != nil {
		fmt.Println("Encountered an error while unmarshalling the data:", err)
		return err
	}
	fmt.Println(mLocation)

	areas := mLocation.Results
	for i := 0; i < len(areas); i++ {
		fmt.Println(areas[i].Name)
	}

	if mLocation.NextUrl == nil {
		conf.NextUrl = "https://pokeapi.co/api/v2/location-area/?offset=20&limit=20"
	} else {
		conf.NextUrl = mLocation.NextUrl
	}
	conf.PreviousUrl = *mLocation.PreviousUrl
	return nil
}

func main() {
	fmt.Println("Welcome to the Pokedex!")
	initCommands()
	scanner := bufio.NewScanner(os.Stdin)

	var conf config
	conf.NextUrl= "https://pokeapi.co/api/v2/location-area/?offset=20&limit=20"


	commandHelp(&conf)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := cleanInput(scanner.Text())
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
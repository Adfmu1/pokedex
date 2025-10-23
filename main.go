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
			description:	"Displays 20 names of location areas, call again to go forward",
			callback:		commandMap,
		},
		"mapb": {
			name:			"mapb",
			description:	"Displays 20 previous names of location areas, call again to go back",
			callback:		commandMapb,
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
	// get url
	url := conf.NextUrl
	if url == "" {
		fmt.Println("No more pages forward, please use command 'mapb'")
		return nil
	}

	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Encountered an error while trying to get url:", err)
		return err
	}
	defer res.Body.Close()
	// body
	body, err := io.ReadAll(res.Body)
	if res.StatusCode > 299 {
		fmt.Println("Bad status code:", res.StatusCode)
		return err
	}
	if err != nil {
		fmt.Println("Encountered an error while retrieving body:", err)
		return err
	}
	// unmarshal the data
	mLocation := LocationAreas{}
	err = json.Unmarshal(body, &mLocation)
	if err != nil {
		fmt.Println("Encountered an error while unmarshalling the data:", err)
		return err
	}
	// print locations
    for _, a := range mLocation.Results {
        fmt.Println(a.Name)
    }
	// set new previous and next URLs
	if mLocation.NextUrl == "" {
		conf.NextUrl = ""
		fmt.Println("No more pages forward, please use command 'mapb'")
	} else {
		conf.NextUrl = mLocation.NextUrl
	}
	if mLocation.PreviousUrl != nil {
		conf.PreviousUrl = *mLocation.PreviousUrl
	} else {
		conf.PreviousUrl = ""
	}
	return nil
}

func commandMapb(conf *config) error {
    // get url
    url := conf.PreviousUrl
    if url == "" {
        fmt.Println("No more pages backward, please use command 'map'")
        return nil
    }

    res, err := http.Get(url)
    if err != nil {
        fmt.Println("Encountered an error while trying to get url:", err)
        return err
    }
    defer res.Body.Close()

    // read body
    body, err := io.ReadAll(res.Body)
    if err != nil {
        fmt.Println("Encountered an error while retrieving body:", err)
        return err
    }
    if res.StatusCode < 200 || res.StatusCode > 299 {
        return fmt.Errorf("bad status %d: %s", res.StatusCode, string(body))
    }

    // unmarshal the data
    mLocation := LocationAreas{}
    if err := json.Unmarshal(body, &mLocation); err != nil {
        fmt.Println("Encountered an error while unmarshalling the data:", err)
        return err
    }

    // print locations
    for _, a := range mLocation.Results {
        fmt.Println(a.Name)
    }

    // set new previous and next URLs
    conf.NextUrl = mLocation.NextUrl

    if mLocation.PreviousUrl == nil {
        conf.PreviousUrl = ""
        fmt.Println("No more pages backward, please use command 'map'")
    } else {
        conf.PreviousUrl = *mLocation.PreviousUrl
    }

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
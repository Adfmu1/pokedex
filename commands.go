package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func initCommands() {
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
			description: "Displays 20 names of location areas, call again to go forward",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays 20 previous names of location areas, call again to go back",
			callback:    commandMapb,
		},
	}
}

func cleanInput(text string) []string {
	str := strings.Fields(text)
	for i := range str {
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
	fmt.Println("Usage:")
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
	// check if data is in cache already
	data, found := conf.Cache.Get(url)
	// if not make a request to api
	if !found {
		// get api data
		res, err := http.Get(url)
		if err != nil {
			fmt.Println("Encountered an error while trying to get url:", err)
			return err
		}
		defer res.Body.Close()
		// read body
		data, err = io.ReadAll(res.Body)
		if res.StatusCode > 299 {
			fmt.Println("Bad status code:", res.StatusCode)
			return err
		}
		if err != nil {
			fmt.Println("Encountered an error while retrieving body:", err)
			return err
		}
	}
	// unmarshal the data
	mLocation := LocationAreas{}
	err := json.Unmarshal(data, &mLocation)
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
	// if data is correct not found in cache save to cache
	if !found {
		conf.Cache.Add(url, data)
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
	// check if data is in cache already
	data, found := conf.Cache.Get(url)
	// if not make a request to api
	if !found {
		// get api data
		res, err := http.Get(url)
		if err != nil {
			fmt.Println("Encountered an error while trying to get url:", err)
			return err
		}
		defer res.Body.Close()
		// read body
		data, err = io.ReadAll(res.Body)
		if res.StatusCode > 299 {
			fmt.Println("Bad status code:", res.StatusCode)
			return err
		}
		if err != nil {
			fmt.Println("Encountered an error while retrieving body:", err)
			return err
		}
	}
	// unmarshal the data
	mLocation := LocationAreas{}
	if err := json.Unmarshal(data, &mLocation); err != nil {
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
	// if data is correct not found in cache save to cache
	if !found {
		conf.Cache.Add(url, data)
	}
	return nil
}

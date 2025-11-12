package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
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
		"explore": {
			name:        "explore",
			description: "Displays list of all the Pokémon located in given location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Throws Pokeball at Pokemon with given name",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Prints name, height, weight, stats and types of Pokemon (if already caught)",
			callback:    commandInspect,
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

func commandExit(params ...any) error {
	defer os.Exit(0)
	fmt.Println("Closing the Pokedex... Goodbye!")
	return nil
}

func commandHelp(params ...any) error {
	fmt.Println("Usage:")
	for _, comm := range commands {
		fmt.Println(comm.name, ":", comm.description)
	}
	return nil
}

func commandMap(params ...any) error {
	// check if parameter is of *conf type
	if len(params) != 1 {
		return errors.New("wrong amount of parameter, should be 1")
	}
	if _, ok := params[0].(*config); !ok {
		return errors.New("wrong type of parameter, should be pointer to config type struct")
	}
	conf := params[0].(*config)
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

func commandMapb(params ...any) error {
	// check if parameter is of *conf type
	if len(params) != 1 {
		return errors.New("wrong amount of parameter, should be 1")
	}
	if _, ok := params[0].(*config); !ok {
		return errors.New("wrong type of parameter, should be pointer to config type struct")
	}
	conf := params[0].(*config)
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

// should accept *config and a string (location name)
func commandExplore(params ...any) error {
	// check if parameter is of *conf type
	if len(params) != 2 {
		return errors.New("wrong amount of parameters, should be 2")
	}
	if _, ok := params[0].(*config); !ok {
		return errors.New("wrong type of parameter[0], should be pointer to config type struct")
	} else if _, ok := params[1].(string); !ok {
		return errors.New("wrong type of parameter[1], should be string")
	}
	conf := params[0].(*config)
	location := params[1].(string)

	if len(location) == 0 {
		return errors.New("no location given")
	}
	// create full url
	const urlStart = "https://pokeapi.co/api/v2/location-area/"
	url := urlStart + location
	// check if data in cache
	data, found := conf.Cache.Get(url)
	// if not make request to api
	if !found {
		// make call to api
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
	// create struct just for poke
	var pokedata struct {
		PokemonEncounters []struct {
			Pokemon struct {
				Name string `json:"name"`
			} `json:"pokemon"`
		} `json:"pokemon_encounters"`
	}
	// unmarshal the data
	err := json.Unmarshal(data, &pokedata)
	if err != nil {
		fmt.Println("Encountered an error while unmarshalling the data:", err)
		return err
	}
	// print pokemon names
	for _, encounter := range pokedata.PokemonEncounters {
		fmt.Println(encounter.Pokemon.Name)
	}
	// if not in cache, save
	if !found {
		conf.Cache.Add(url, data)
	}
	return nil
}

// should accept a string (pokemon name) and a Pokedex struct
func commandCatch(params ...any) error {
	// check if correct number of parameters
	if len(params) != 2 {
		return errors.New("wrong number of parameters, should be 2")
	}
	// check if correct type of parameter
	pokemonName, ok := params[0].(string)
	if !ok {
		return errors.New("wrong type of parameter, should be string")
	}
	pokedex, ok := params[1].(*Pokedex)
	if !ok {
		return errors.New("wrong type of parameter, should be pntr to a pokedex struct")
	}
	pokemonName = strings.ToLower(pokemonName)
	// check if pokemon is in pokedex
	if _, ok := pokedex.Pokemons[pokemonName]; ok {
		fmt.Printf("Pokemon %s has been already caught\n", pokemonName)
		return nil
	}
	// get URL
	url := "https://pokeapi.co/api/v2/pokemon/"
	fullUrl := url + pokemonName
	res, err := http.Get(fullUrl)
	if err != nil {
		fmt.Println("Encountered an error while trying to get url:", err)
		return err
	}
	defer res.Body.Close()
	// read body
	data, err := io.ReadAll(res.Body)
	if res.StatusCode > 299 {
		fmt.Println("Bad status code:", res.StatusCode)
		return err
	}
	if err != nil {
		fmt.Println("Encountered an error while retrieving body:", err)
		return err
	}
	// unmarshal the data
	pokemon := Pokemon{}
	err = json.Unmarshal(data, &pokemon)
	if err != nil {
		fmt.Println("Encountered an error while unmarshalling the data:", err)
		return err
	}
	// try to catch a pokemon
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	if chance := rand.Intn(1000); chance <= pokemon.BaseExp*2 {
		fmt.Printf("%s escaped!\n", pokemonName)
		return nil
	}
	fmt.Printf("%s was caught!\n", pokemonName)
	pokedex.Pokemons[pokemon.Name] = pokemon
	return nil
}

// should accept a string (pokemon name) and a Pokedex struct
func commandInspect(params ...any) error {
	// check if correct number of parameters
	if len(params) != 2 {
		return errors.New("wrong number of parameters, should be 2")
	}
	// check if correct type of parameter
	pokemonName, ok := params[0].(string)
	if !ok {
		return errors.New("wrong type of parameter, should be string")
	}
	pokedex, ok := params[1].(*Pokedex)
	if !ok {
		return errors.New("wrong type of parameter, should be pntr to a pokedex struct")
	}
	pokemonName = strings.ToLower(pokemonName)
	// check if pokemon is in pokedex
	pokemon, ok := pokedex.Pokemons[pokemonName]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for i := 0; i < len(pokemon.Stats); i++ {
		fmt.Printf("\t-%s: %d\n", pokemon.Stats[i].Stat.Name, pokemon.Stats[i].BaseValue)
	}
	fmt.Println("Types:")
	for i := 0; i < len(pokemon.Types); i++ {
		fmt.Printf("\t-%s\n", pokemon.Types[i].Type.NamedRes)
	}
	return nil
}

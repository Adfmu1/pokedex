package main

import (
	"github.com/Adfmu1/pokedex/internal/pokecache"
)

// also serves as cache
type config struct {
	NextUrl     string
	PreviousUrl string
	Cache       *pokecache.Cache
}

// ============ POKEMONS ============
type Pokedex struct {
	Pokemons map[string]Pokemon
}

type Pokemon struct {
	Name   string     `json:"name"`
	Height int        `json:"height"`
	Weight int        `json:"weight"`
	Stats  []PokeStat `json:"stats"`
	Types  []PokeType `json:"types"`
}

type PokeStat struct {
	Stat struct {
		Name string `json:"name"`
	} `json:"stat"`
	BaseValue int `json:"base_stat"`
}

type PokeType struct {
	Type struct {
		NamedRes string `json:"name"`
	} `json:"type"`
}

// ============ LOCATIONS ============
type LocationAreas struct {
	Results     []Area  `json:"results"`
	NextUrl     string  `json:"next"`
	PreviousUrl *string `json:"previous"`
}

type Area struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

// ============ COMMANDS ============
type cliCommand struct {
	name        string
	description string
	callback    func(params ...any) error
}

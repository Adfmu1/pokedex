package main

import (
	"github.com/Adfmu1/pokedex/internal/pokecache"
)

type config struct {
	NextUrl     string
	PreviousUrl string
	Cache       *pokecache.Cache
}

type LocationAreas struct {
	Results     []Area  `json:"results"`
	NextUrl     string  `json:"next"`
	PreviousUrl *string `json:"previous"`
}

type Area struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

type cliCommand struct {
	name        string
	description string
	callback    func(params ...any) error
}

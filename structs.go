package main

type config struct {
	nextUrl			string
	previousUrl		string
}

type LocationAreas struct {
	Results			[]Area	`json:"results"`
	nextUrl			string	`json:"next"`
	previousUrl		string	`json:"previous"`
}

type Area struct {
	ID   			int    `json:"id"`
    Name 			string `json:"name"`
}
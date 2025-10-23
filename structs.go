package main

type config struct {
	NextUrl			string
	PreviousUrl		string
}

type LocationAreas struct {
	Results			[]Area	`json:"results"`
	NextUrl			string	`json:"next"`
	PreviousUrl		*string	`json:"previous"`
}

type Area struct {
	Url   			string	`json:"url"`
    Name 			string	`json:"name"`
}
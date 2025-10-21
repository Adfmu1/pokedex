package main

import (
	"fmt"
	"strings"
	)

func cleanInput(text string) []string {
	str := strings.Fields(text)
	for i := 0; i < len(str); i++ {
		str[i] = strings.ToLower(str[i])
	}

	return str
}

func main() {
	fmt.Println("Hello, World!")
}
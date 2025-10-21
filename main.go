package main

import (
	"bufio"
	"os"
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
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := cleanInput(scanner.Text())
		
		fmt.Println("Your command was:", input[0])

	}
}
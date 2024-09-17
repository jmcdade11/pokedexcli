package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func startRepl(config *config) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("Pokedex > ")
		scanner.Scan()
		words := cleanInput(scanner.Text())
		if len(words) == 0 {
			continue
		}
		commandString := words[0]
		args := []string{}
		if len(words) > 1 {
			args = words[1:]
		}
		if command, ok := getCommands()[commandString]; ok {
			err := command.callback(config, args...)
			if err != nil {
				fmt.Println(err)
			}
			continue
		} else {
			fmt.Println("Unknown command")
			continue
		}
	}
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

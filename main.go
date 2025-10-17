package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"io"
	"net/http"
	"encoding/json"
)

type cfg struct {
	Next        *string
	Previous    *string
	BaseURL     string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*cfg) error
}

type mapPage struct {
	Count       int `json:"count"`
	Next        *string `json:"next"`
	Previous        *string `json:"previous"`
	Results     []locationArea `json:"results"`
}

type locationArea struct {
	Name        string `json:"name"`
	Url         string `json:"url"`
}

var commands map[string]cliCommand

func cleanInput(text string) []string{
	lowerText := strings.ToLower(text)
	lowerText = strings.TrimSpace(lowerText)
	return strings.Fields(lowerText)
}

func commandExit(c *cfg) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *cfg) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n")
	for _, val := range(commands){
		fmt.Printf("%v: %v\n", val.name, val.description)
	}
	return nil
}

func commandMap(c *cfg) error {
	var url string
	if c.Next != nil {
		url = *c.Next
	} else{
		url = c.BaseURL
	}
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode > 299 {
		return fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil{
		return err
	}
	var data mapPage
	err = json.Unmarshal(body, &data)
	c.Next = data.Next
	c.Previous = data.Previous
	if err != nil {
		return err
	}
	for _, loc := range(data.Results){
		fmt.Println(loc.Name)
	}
	return nil
}

func commandMapb(c *cfg) error {
	var url string
	if c.Previous != nil {
		url = *c.Previous
	} else{
		fmt.Println("you're on the first page")
		return nil
	}
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode > 299 {
		return fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil{
		return err
	}
	var data mapPage
	err = json.Unmarshal(body, &data)
	c.Next = data.Next
	if data.Previous != nil{
		c.Previous = data.Previous
	} else {
		c.Previous = nil
	}
	for _, loc := range(data.Results){
		fmt.Println(loc.Name)
	}
	return nil
}

func init(){
	commands = map[string]cliCommand{
	"help": {
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	},
    "exit": {
        name:        "exit",
        description: "Exit the Pokedex",
        callback:    commandExit,
    },
	"map": {
		name:        "map",
		description: "Lists locations in the Pokemon world (or go to the next page)",
		callback:    commandMap,
	},
	"mapb": {
		name:        "mapb",
		description: "Goes to the previous page of map locations",
		callback:    commandMapb,
	},
}
}

func NewConfig() *cfg {
	base := "https://pokeapi.co/api/v2/location-area/"
	res := &cfg{
		Next:     &base,
		Previous: nil,
		BaseURL:  base,
	}
	return res
}

func main(){
	config := NewConfig()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			prompt := scanner.Text()
			command := cleanInput(prompt)
			if len(command) >= 1{
				_, ok := commands[command[0]]
				if ok {
					commands[command[0]].callback(config)
				}else{
					fmt.Println("Unknown command")
				}
			}
		}else{
			fmt.Println(scanner.Err())
			break
		}
	}
}
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	rcon "github.com/gtaylor/factorio-rcon"
)

var (
	rconServerAddress  = os.Getenv("RCON_SERVER_ADDRESS")
	rconServerPassword = os.Getenv("RCON_SERVER_PASSWORD")
)

func main() {
	r, err := rcon.Dial(rconServerAddress)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	err = r.Authenticate(rconServerPassword)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		trimmedText := strings.TrimSpace(text)

		response, err := r.Execute(trimmedText)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Response: %+v\n", response.Body)
	}
}

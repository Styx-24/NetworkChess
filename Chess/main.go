package main

import (
	"TP/client"
	"TP/server"
	"fmt"
	"os"
)

func main() {

	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Indiquez un option de lancement")
		return
	}

	if arguments[1] == "s" {
		server.Server()
	} else if arguments[1] == "c" {
		client.Client()
	}

}

package main

import "fmt"

func printHelp() {
	fmt.Print("Usage: dish [FLAGS] SOURCE\n\n")
	fmt.Print("A lightweight, one-shot socket checker\n\n")
	fmt.Println("SOURCE must be a file path leading to a JSON file with a list of sockets to be checked or a URL leading to a remote JSON API from which the list of sockets can be retrieved")
	fmt.Println("Use the `-h` flag for a list of available flags")
}

package main

import (
	"log"
	"os"

	"github.com/ahmedmkamal/go-args"
)

func main() {
	parser := args.NewParser(os.Args[1:])
	if err := parser.Parse(); err != nil {
		log.Fatal(err)
	}

	if parser.HasOption("-h") {
		log.Println("Help Script!")
		return
	}

	cmd, exist := parser.At(0)
	if !exist {
		log.Fatal("no positional arg found at index 0")
	}
	log.Printf("positional arg found at index 0 = %s\n", cmd)
}

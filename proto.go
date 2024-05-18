package main

import "fmt"

type Command struct {
	// ??
}

func parseCommand(rawMsg []byte) (Command, error) {
	t := rawMsg[0]
	fmt.Println("received message", string(rawMsg))
	switch t {
	case '*':
		fmt.Println(rawMsg[1:])
	}
	return Command{}, nil
}

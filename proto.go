package main

import (
	"bytes"
	"fmt"
	"github.com/tidwall/resp"
	"io"
	"log"
)

const (
	CommandSet = "SET"
)

type Command interface {
	//
}

type SetCommand struct {
	key, value string
}

func parseCommand(raw string) (Command, error) {
	rd := resp.NewReader(bytes.NewBufferString(raw))
	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if v.Type() == resp.Array {
			for _, value := range v.Array() {
				switch value.String() {
				case CommandSet:
					if len(v.Array()) != 3 {
						log.Fatal("invalid command")
					}
					cmd := SetCommand{
						key:   v.Array()[1].String(),
						value: v.Array()[2].String(),
					}
					return cmd, nil
				default:
				}
			}
		}
		return nil, fmt.Errorf("invalid or unknown command recieved: %s", raw)
	}
	return nil, fmt.Errorf("invalid or unknown command recieved: %s", raw)
}

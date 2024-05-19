package main

import (
	"bytes"
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
			for _, v := range v.Array() {
				switch v.String() {
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
	}
	return "", nil
}

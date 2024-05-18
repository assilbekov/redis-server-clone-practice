package main

import (
	"bytes"
	"fmt"
	"github.com/tidwall/resp"
	"io"
	"log"
	"testing"
)

func TestProtocol(t *testing.T) {
	raw := "*3\r\n$3\r\nset\r\n$6\r\nleader\r\n$7\r\nCharlie\r\n"
	raw += "*3\r\n$3\r\nset\r\n$8\r\nfollower\r\n$6\r\nSkyler\r\n"
	rd := resp.NewReader(bytes.NewBufferString(raw))
	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Read %s\n", v.Type())
		if v.Type() == resp.Array {
			for i, v := range v.Array() {
				fmt.Printf("  #%d %s, value: '%s'\n", i, v.Type(), v)
			}
		}
	}
}

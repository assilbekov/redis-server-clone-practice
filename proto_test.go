package main

import (
	"bytes"
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
		if v.Type() == resp.Array {
			for _, v := range v.Array() {
				switch v.String() {
				case CommandSet:
				default:
				}
			}
		}
	}
}

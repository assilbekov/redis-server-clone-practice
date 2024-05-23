package client

import (
	"bytes"
	"context"
	"github.com/tidwall/resp"
	"io"
	"log"
	"net"
)

type Client struct {
	addr string
	conn net.Conn
}

func NewClient(addr string) *Client {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	return &Client{addr: addr, conn: conn}
}

func (c *Client) Set(ctx context.Context, key, value string) error {
	buf := &bytes.Buffer{}
	wr := resp.NewWriter(buf)
	wr.WriteArray([]resp.Value{
		resp.StringValue("SET"),
		resp.StringValue(key),
		resp.StringValue(value),
	})

	_, err := io.Copy(c.conn, buf)
	return err
}

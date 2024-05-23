package client

import (
	"bytes"
	"context"
	"github.com/tidwall/resp"
	"io"
	"net"
)

type Client struct {
	addr string
	conn net.Conn
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Set(ctx context.Context, key, value string) error {
	if c.conn == nil {
		conn, err := net.Dial("tcp", "localhost:5001")
		if err != nil {
			return err
		}
		c.conn = conn
	}

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

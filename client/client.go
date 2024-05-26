package client

import (
	"bytes"
	"context"
	"fmt"
	"github.com/tidwall/resp"
	"net"
)

type Client struct {
	addr string
	conn net.Conn
}

func NewClient(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %v", err)
	}

	return &Client{addr: addr, conn: conn}, nil
}

func (c *Client) Set(ctx context.Context, key, value string) error {
	buf := &bytes.Buffer{}
	wr := resp.NewWriter(buf)
	wr.WriteArray([]resp.Value{
		resp.StringValue("SET"),
		resp.StringValue(key),
		resp.StringValue(value),
	})

	_, err := c.conn.Write(buf.Bytes())
	return err
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	buf := &bytes.Buffer{}
	wr := resp.NewWriter(buf)
	wr.WriteArray([]resp.Value{
		resp.StringValue("GET"),
		resp.StringValue(key),
	})

	_, err := c.conn.Write(buf.Bytes())
	if err != nil {
		return "", err
	}

	b := make([]byte, 1024)
	n, err := c.conn.Read(b)
	return string(b[:n]), err
}

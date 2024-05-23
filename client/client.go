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
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Set(ctx context.Context, key, value string) error {
	conn, err := net.Dial("tcp", "localhost:5001")
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	wr := resp.NewWriter(buf)
	wr.WriteArray([]resp.Value{
		resp.StringValue("SET"),
		resp.StringValue(key),
		resp.StringValue(value),
	})

	_, err = io.Copy(conn, buf)
	return err
}

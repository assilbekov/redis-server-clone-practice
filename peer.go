package main

import (
	"log/slog"
	"net"
)

type Peer struct {
	conn net.Conn
}

func NewPeer(conn net.Conn) *Peer {
	return &Peer{conn: conn}
}

func (p *Peer) readLoop() error {
	buf := make([]byte, 1024)
	for {
		n, err := p.conn.Read(buf)
		if err != nil {
			slog.Error("failed to read from connection", "err", err)
			return err
		}
		msgBuf := make([]byte, n)
		copy(msgBuf, buf[:n])
	}
}

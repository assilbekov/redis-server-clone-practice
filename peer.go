package main

import (
	"fmt"
	"github.com/tidwall/resp"
	"io"
	"log"
	"net"
)

type Peer struct {
	conn  net.Conn
	msgCh chan Message
	delCh chan *Peer
}

func (p *Peer) Send(msg []byte) (int, error) {
	return p.conn.Write(msg)
}

func NewPeer(conn net.Conn, msgCh chan Message, delCh chan *Peer) *Peer {
	return &Peer{conn: conn, msgCh: msgCh, delCh: delCh}
}

func (p *Peer) readLoop() error {
	rd := resp.NewReader(p.conn)

	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			p.delCh <- p
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if v.Type() == resp.Array {
			for _, value := range v.Array() {
				var cmd Command
				switch value.String() {
				case CommandGet:
					if len(v.Array()) != 2 {
						return fmt.Errorf("invalid get command")
					}
					cmd = GetCommand{
						key: v.Array()[1].Bytes(),
					}
				case CommandSet:
					if len(v.Array()) != 3 {
						return fmt.Errorf("invalid set command")
					}
					cmd = SetCommand{
						key:   v.Array()[1].Bytes(),
						value: v.Array()[2].Bytes(),
					}
				case CommandHello:
					if len(v.Array()) != 2 {
						return fmt.Errorf("invalid hello command")
					}
					cmd = HelloCommand{
						value: v.Array()[1].String(),
					}
				default:
					panic("unknown command: " + value.String())
				}

				p.msgCh <- Message{peer: p, cmd: cmd}
			}
		}
	}
	return nil
}

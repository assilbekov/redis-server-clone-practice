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

		var cmd Command
		if v.Type() == resp.Array {
			rawCmd := v.Array()[0]
			fmt.Println("rawCmd =>", rawCmd.String())
			fmt.Println("this should be a command", rawCmd.String())

			switch rawCmd.String() {
			case CommandGet:
				cmd = GetCommand{
					key: v.Array()[1].Bytes(),
				}
			case CommandSet:
				cmd = SetCommand{
					key:   v.Array()[1].Bytes(),
					value: v.Array()[2].Bytes(),
				}
			case CommandHello:
				cmd = HelloCommand{
					value: v.Array()[1].String(),
				}
			default:
				fmt.Println("invalid command", rawCmd)
			}
			p.msgCh <- Message{peer: p, cmd: cmd}
		}
	}
	return nil
}

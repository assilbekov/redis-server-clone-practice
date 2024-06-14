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
			case CommandSet:
			case CommandHello:
				cmd = HelloCommand{
					value: v.Array()[1].String(),
				}
			}
			p.msgCh <- Message{peer: p, cmd: cmd}
			/*for _, value := range v.Array() {
				fmt.Println("value =>", value.String())
				//var rawCmd Command
				switch value.String() {
				case CommandGet:
					if len(v.Array()) != 2 {
						return fmt.Errorf("invalid get command")
					}
					rawCmd := GetCommand{
						key: v.Array()[1].Bytes(),
					}

					p.msgCh <- Message{peer: p, rawCmd: rawCmd}
				case CommandSet:
					if len(v.Array()) != 3 {
						return fmt.Errorf("invalid set command")
					}
					rawCmd := SetCommand{
						key:   v.Array()[1].Bytes(),
						value: v.Array()[2].Bytes(),
					}
					p.msgCh <- Message{peer: p, rawCmd: rawCmd}
				case CommandHello:
					rawCmd := HelloCommand{
						value: v.Array()[1].String(),
					}
					p.msgCh <- Message{peer: p, rawCmd: rawCmd}
				}

			}*/
		}
	}
	return nil
}

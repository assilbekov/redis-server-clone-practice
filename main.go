package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"redis-server-clone-practice/client"
	"time"
)

const defaultListenAddr = ":5001"

type Config struct {
	ListenAddr string
}

type Message struct {
	data []byte
	peer *Peer
}

type Server struct {
	Config
	peers     map[*Peer]bool
	ln        net.Listener
	addPeerCh chan *Peer
	quitCh    chan struct{}
	msgCh     chan Message

	//
	kv *KV
}

func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = defaultListenAddr
	}
	return &Server{
		Config:    cfg,
		peers:     make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
		quitCh:    make(chan struct{}),
		msgCh:     make(chan Message),
		kv:        NewKV(),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}

	s.ln = ln

	go s.loop()

	slog.Info("server started", "addr", s.ListenAddr)

	return s.acceptLoop()
}

func (s *Server) handleMessage(msg Message) error {
	cmd, err := parseCommand(string(msg.data))
	if err != nil {
		return err
	}
	switch v := cmd.(type) {
	case SetCommand:
		return s.kv.Set(v.key, v.value)
	case GetCommand:
		value, ok := s.kv.Get(v.key)
		if !ok {
			return fmt.Errorf("key not found: %s", v.key)
		}
		_, err := msg.peer.Send(value)
		if err != nil {
			slog.Error("failed to send message", "err", err)
		}
	}

	return nil
}

func (s *Server) loop() {
	for {
		select {
		case msg := <-s.msgCh:
			if err := s.handleMessage(msg); err != nil {
				slog.Error("failed to handle message", "err", err)
			}
			//fmt.Println("received message", string(rawMsg))
		case <-s.quitCh:
			return
		case p := <-s.addPeerCh:
			s.peers[p] = true
		}
	}
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("failed to accept connection", "err", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	peer := NewPeer(conn, s.msgCh)
	s.addPeerCh <- peer
	if err := peer.readLoop(); err != nil {
		slog.Info("peer disconnected", "remoteAddr", conn.RemoteAddr())
	}
}

func main() {
	server := NewServer(Config{})
	go func() {
		log.Fatal(server.Start())
	}()

	time.Sleep(time.Second * 2)

	c, err := client.NewClient("localhost:5001")
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		if err := c.Set(
			context.Background(),
			fmt.Sprintf("leader_%d", i),
			fmt.Sprintf("Charlie_%d", i),
		); err != nil {
			log.Fatal(err)
		}

		time.Sleep(time.Millisecond * 100)

		if v, err := c.Get(context.Background(), fmt.Sprintf("leader_%d", i)); err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(v)
		}
	}

	time.Sleep(time.Second * 2)
	fmt.Println(server.kv.data)
	// select {} // we are blocking here to keep the server running
}

package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
)

const defaultListenAddr = ":5001"

type Config struct {
	ListenAddr string
}

type Message struct {
	cmd  Command
	peer *Peer
}

type Server struct {
	Config
	peers     map[*Peer]bool
	ln        net.Listener
	addPeerCh chan *Peer
	delPeerCh chan *Peer
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
		delPeerCh: make(chan *Peer),
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

	slog.Info("redis server started", "addr", s.ListenAddr)

	return s.acceptLoop()
}

func (s *Server) handleMessage(msg Message) error {
	switch v := msg.cmd.(type) {
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
			slog.Info("new peer connected", "remoteAddr", p.conn.RemoteAddr())
			s.peers[p] = true
		case p := <-s.delPeerCh:
			slog.Info("peer disconnected", "remoteAddr", p.conn.RemoteAddr())
			delete(s.peers, p)
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
	listenAddr := flag.String("listenAddr", defaultListenAddr, "server listen address")
	flag.Parse()
	server := NewServer(Config{
		ListenAddr: *listenAddr,
	})
	log.Fatal(server.Start())
}

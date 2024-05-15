package main

import (
	"log"
	"log/slog"
	"net"
)

const defaultListenAddr = ":5001"

type Config struct {
	ListenAddr string
}

type Server struct {
	Config
	peers     map[*Peer]bool
	ln        net.Listener
	addPeerCh chan *Peer
	quitCh    chan struct{}
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

func (s *Server) loop() {
	for {
		select {
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
	peer := NewPeer(conn)
	s.addPeerCh <- peer

	slog.Info("new peer connected", "remoteAddr", conn.RemoteAddr())
	peer.readLoop()
}

func main() {
	server := NewServer(Config{})
	log.Fatal(server.Start())
}

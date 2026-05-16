package server

import (
	"io"
	"log/slog"
	"net"
	"sync"
)

type Server struct {
	listener net.Listener
	peers    *Peers
}

type Peers struct {
	mu    sync.Mutex
	conns map[net.Conn]struct{}
}

func NewServer(port string) (*Server, error) {

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		slog.Error("Failed at New server", "err", err)
		return nil, err
	}

	return &Server{
		listener: ln,
		peers: &Peers{
			conns: make(map[net.Conn]struct{}),
		},
	}, nil
}
func (s *Server) Close() error {
	return s.listener.Close()
}

func (s *Server) Serve() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			slog.Error("accept error", "error", err)
			continue
		}

		s.addPeer(conn)
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		s.removePeer(conn)
		conn.Close()
	}()

	addr := conn.RemoteAddr().String()
	slog.Info("client connected", "addr", addr)

	buffer := make([]byte, 4096)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				slog.Info("client disconnected", "addr", addr)
			} else {
				slog.Error("read error", "error", err)
			}
			return
		}

		data := buffer[:n]
		slog.Info("received", "bytes", n, "data", data)

		cmd, err := readCommand(data)
		if err != nil {
			respError(conn, err.Error())
			continue
		}
		resp(conn, cmd)
		// if _, err := conn.Write([]byte("+OK\r\n")); err != nil {
		// slog.Error("write error", "error", err)
		// 	return
		// }
	}
}

func (s *Server) addPeer(conn net.Conn) {
	s.peers.mu.Lock()
	defer s.peers.mu.Unlock()
	s.peers.conns[conn] = struct{}{}
}

func (s *Server) removePeer(conn net.Conn) {
	s.peers.mu.Lock()
	defer s.peers.mu.Unlock()
	delete(s.peers.conns, conn) // i dont like this
}

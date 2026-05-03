package server

import (
	"bufio"
	"net"
	"testing"
	"time"
)

func startTestServer(t *testing.T) *Server {
	srv, err := NewServer(":6575")
	if err != nil {
		t.Fatalf("failed to start server: %v", err)
	}

	go srv.Serve()

	return srv
}

func TestServer_ResponseOK(t *testing.T) {
	srv := startTestServer(t)
	defer srv.Close()

	addr := srv.listener.Addr().String()

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte("PING\n"))
	if err != nil {
		t.Fatalf("write failed: %v", err)
	}

	reader := bufio.NewReader(conn)
	resp, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}

	expected := "+OK\r\n"
	if resp != expected {
		t.Fatalf("expected %q, got %q", expected, resp)
	}
}

func TestServer_MultipleConnections(t *testing.T) {
	srv := startTestServer(t)
	defer srv.Close()

	addr := srv.listener.Addr().String()

	for i := 0; i < 5; i++ {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			t.Fatalf("connection %d failed: %v", i, err)
		}

		go func(c net.Conn) {
			defer c.Close()
			c.Write([]byte("hello\n"))

			buf := make([]byte, 16)
			c.SetReadDeadline(time.Now().Add(time.Second))
			_, err := c.Read(buf)
			if err != nil {
				t.Errorf("read error: %v", err)
			}
		}(conn)
	}
}

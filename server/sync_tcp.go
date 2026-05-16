package server

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/abdooman21/personalredis/core"
)

// seperate server logic from core logic
func readCommand(buf []byte) (*core.Command, error) {

	return core.DecodeRESP(buf)

}
func respError(conn net.Conn, err string) {
	if _, err := conn.Write([]byte(fmt.Sprintf("-%s\r\n", err))); err != nil {
		slog.Error("write error", "error", err)
		return
	}
}

func resp(conn net.Conn, cmd *core.Command) {
	switch cmd.Cmd[0] {
	case 'P':
		if cmd.Args[0] == "" {
			if _, err := conn.Write([]byte("PONG\r\n")); err != nil {
				slog.Error("write error", "error", err)
			}
			return
		}
		if _, err := conn.Write([]byte(fmt.Sprintf("\"%s\"", cmd.Args[0]))); err != nil {
			slog.Error("write error", "error", err)
			return
		}
		return

	}
}

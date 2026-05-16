package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/abdooman21/personalredis/server"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load(".env")

	portString := os.Getenv("PORT")

	srv, err := server.NewServer(portString)
	if err != nil {
		os.Exit(1)
	}
	defer srv.Close()

	fmt.Println("Custom Redis server listening on port ...", portString)
	go srv.Serve()

	// for {
	// 	conn, err := srv.Listen.Accept()
	// 	if err != nil {
	// 		log.Println("Accept error:", err)
	// 		continue
	// 	}

	// 	go handleConnection(conn)
	// }

	sigs := make(chan os.Signal, 1)

	// 2. Register the channel to receive specific signals.
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs
	fmt.Printf("\nReceived signal: %s. Cleaning up...\n", sig)
}

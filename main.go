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

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs
	fmt.Printf("\nReceived signal: %s. Cleaning up...\n", sig)
}

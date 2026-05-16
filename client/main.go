package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	p := flag.Int("p", 6575, "number of workers")
	flag.Parse()

	//127.0.0.1
	addr := fmt.Sprintf("127.0.0.1:%v", *p)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("redis> ")

		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("read error:", err)
			return
		}

		// send to redis
		_, err = conn.Write([]byte(text))
		if err != nil {
			fmt.Println("write error:", err)
			return
		}

		// read response
		buffer := make([]byte, 4096)

		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("server closed:", err)
			return
		}

		fmt.Println("response:", string(buffer[:n]))
	}

}

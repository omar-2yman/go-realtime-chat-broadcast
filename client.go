package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
)

type Message struct {
	From    int
	Content string
}

type ClientRPC struct{}

func (c *ClientRPC) Receive(msg Message, _ *bool) error {
	fmt.Printf("\n[User %d] %s\n", msg.From, msg.Content)
	return nil
}

func main() {
	// connect to server
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		log.Fatal(err)
	}
	server := rpc.NewClient(conn)

	// register local RPC listener
	rpc.Register(new(ClientRPC))

	// start listener on a random port
	ln, _ := net.Listen("tcp", ":0")
	go func() {
		for {
			c, _ := ln.Accept()
			go rpc.ServeConn(c)
		}
	}()

	// get ID
	var myID int
	server.Call("Server.Join", 0, &myID)
	fmt.Println("You joined as User", myID)

	// send address for callbacks
	addr := ln.Addr().String()
	var ignore bool
	server.Call("Server.RegisterClient",
		struct {
			ID   int
			Addr string
		}{myID, addr},
		&ignore,
	)

	// main loop
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Message: ")
		scanner.Scan()
		text := scanner.Text()

		msg := Message{From: myID, Content: text}
		server.Call("Server.Send", msg, &ignore)
	}
}

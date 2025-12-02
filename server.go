package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
)

type Message struct {
	From    int
	Content string
}

type Server struct {
	mu      sync.Mutex
	clients map[int]*rpc.Client
	nextID  int
}

func NewServer() *Server {
	return &Server{
		clients: make(map[int]*rpc.Client),
		nextID:  1,
	}
}

func (s *Server) Join(_ int, reply *int) error {
	s.mu.Lock()
	id := s.nextID
	s.nextID++
	s.mu.Unlock()

	*reply = id

	// broadcast join
	go s.broadcast(Message{
		From:    id,
		Content: fmt.Sprintf("User %d joined", id),
	})

	return nil
}

func (s *Server) RegisterClient(args struct {
	ID   int
	Addr string
}, _ *bool) error {

	conn, err := net.Dial("tcp", args.Addr)
	if err != nil {
		return err
	}

	c := rpc.NewClient(conn)

	s.mu.Lock()
	s.clients[args.ID] = c
	s.mu.Unlock()

	return nil
}

func (s *Server) Send(msg Message, _ *bool) error {
	go s.broadcast(msg)
	return nil
}

func (s *Server) broadcast(msg Message) {
	s.mu.Lock()
	for id, cl := range s.clients {
		if id == msg.From {
			continue
		}
		go cl.Call("ClientRPC.Receive", msg, nil)
	}
	s.mu.Unlock()
}

func main() {
	server := NewServer()
	rpc.Register(server)

	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server listening on 9000...")

	for {
		conn, _ := ln.Accept()
		go rpc.ServeConn(conn)
	}
}

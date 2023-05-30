package main

import (
	"fmt"
	"time"
)

type Server struct {
	exitch chan struct{}
	msgch  chan string
}

func NewServer() *Server {
	return &Server{
		exitch: make(chan struct{}),
		msgch:  make(chan string, 128),
	}
}

func (s *Server) start() {
	fmt.Println("server starting")
	s.loop()
}

func (s *Server) sendMsg(msg string) {
	s.msgch <- msg
}

func (s *Server) loop() {
mainloop:
	for {
		select {
		case <-s.exitch:
			//do some stuff when we need to quit
			fmt.Println("quitting server")
			break mainloop
		case msg := <-s.msgch:
			//do some stuff when we have a message
			s.handleMsg(msg)
		}
	}
	fmt.Println("server shutting down gracefully")
}

func (s *Server) exit() {
	close(s.exitch)
}

func (s *Server) handleMsg(msg string) {
	fmt.Println("we received a message: ", msg)
}

func main() {
	server := NewServer()
	go func() {
		time.Sleep(time.Second * 5)
		server.exit()
	}()
	server.start()
}

package network

import (
	"fmt"
	"time"
)

const timeTicker = time.Second * 5

type (
	ServerOpts struct {
		Transports []Transport
	}

	Server struct {
		ServerOpts

		rpc  chan RPC
		quit chan struct{}
	}
)

func NewServer(opts ServerOpts) *Server {
	return &Server{
		ServerOpts: opts,
		rpc:        make(chan RPC),
		quit:       make(chan struct{}),
	}
}

func (s *Server) Start() {
	s.initTransports()
	ticker := time.NewTicker(timeTicker)

free:
	for {
		select {
		case rpc := <-s.rpc:
			fmt.Printf("From: %v \tPayload: %v \n", rpc.From, string(rpc.Payload))
		case <-s.quit:
			break free
		case <-ticker.C:
			fmt.Printf("do stuff every %v\n", timeTicker)
		}
	}

	fmt.Println("Server shutdown")
}

func (s *Server) initTransports() {
	for idx := range s.Transports {
		go func(tr Transport) {
			for rpc := range tr.Consume() {
				s.rpc <- rpc
			}
		}(s.Transports[idx])
	}
}

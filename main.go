package main

import (
	"time"

	"github.com/isaqueveras/projectx/network"
)

func main() {
	trLocal := network.NewLocalTransport("LOCAL")
	trRemote := network.NewLocalTransport("REMOTE")

	trLocal.Connect(trRemote)
	trRemote.Connect(trLocal)

	go func() {
		for {
			trRemote.SendMessage(trLocal.Addr(), []byte("Ayrton Senna"))
			time.Sleep(time.Second)
		}
	}()

	go func() {
		for {
			trLocal.SendMessage(trRemote.Addr(), []byte("Helley de Abreu Silva Batista"))
			time.Sleep(time.Second * 2)
		}
	}()

	opts := network.ServerOpts{
		Transports: []network.Transport{trLocal, trRemote},
	}

	server := network.NewServer(opts)
	server.Start()
}

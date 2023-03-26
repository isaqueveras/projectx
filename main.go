package main

import (
	"time"

	"github.com/isaqueveras/projectx/network"
	"github.com/isaqueveras/projectx/types"
)

func main() {
	// EP5: https://www.youtube.com/watch?v=kYJyzTkIZjg

	trLocal := network.NewLocalTransport("LOCAL")
	trRemote := network.NewLocalTransport("REMOTE")
	trHome := network.NewLocalTransport("HOME")

	trLocal.Connect(trRemote)
	trRemote.Connect(trLocal)
	trHome.Connect(trLocal)
	trHome.Connect(trRemote)

	go func() {
		for {
			a := types.RandomHash()
			trRemote.SendMessage(trLocal.Addr(), []byte(a.String()))
			time.Sleep(time.Second)
		}
	}()

	go func() {
		for {
			a := types.RandomHash()
			trLocal.SendMessage(trRemote.Addr(), []byte(a.String()))
			time.Sleep(time.Second)
		}
	}()

	go func() {
		for {
			a := types.RandomHash()
			trHome.SendMessage(trLocal.Addr(), []byte(a.String()))
			time.Sleep(time.Second)
		}
	}()

	opts := network.ServerOpts{
		Transports: []network.Transport{trLocal, trRemote, trHome},
		BlockTime:  time.Second,
	}

	network.NewServer(opts).Start()
}

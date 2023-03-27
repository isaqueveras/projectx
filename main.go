package main

import (
	"bytes"
	"math/rand"
	"strconv"
	"time"

	"github.com/isaqueveras/projectx/core"
	"github.com/isaqueveras/projectx/crypto"
	"github.com/isaqueveras/projectx/network"
	"github.com/sirupsen/logrus"
)

func main() {
	trLocal := network.NewLocalTransport("LOCAL")
	trRemote := network.NewLocalTransport("REMOTE")

	trLocal.Connect(trRemote)
	trRemote.Connect(trLocal)

	go func() {
		for {
			time.Sleep(time.Second)
			if err := sendTransaction(trRemote, trLocal.Addr()); err != nil {
				logrus.Error(err)
			}
		}
	}()

	opts := network.ServerOpts{
		Transports: []network.Transport{trLocal},
	}

	s := network.NewServer(opts)
	s.Start()
}

func sendTransaction(tr network.Transport, to network.NetAddr) error {
	tx := core.NewTransaction([]byte(strconv.FormatInt(int64(rand.Intn(1000000000)), 10)))
	tx.Sign(crypto.GeneratePrivateKey())

	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGotTxEncoder(buf)); err != nil {
		return err
	}

	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())
	return tr.SendMessage(to, msg.Bytes())
}

package network

import (
	"fmt"
	"time"

	"github.com/isaqueveras/projectx/core"
	"github.com/isaqueveras/projectx/crypto"
	"github.com/sirupsen/logrus"
)

type (
	ServerOpts struct {
		Transports []Transport
		BlockTime  time.Duration
		PrivateKey *crypto.PrivateKey
	}

	Server struct {
		ServerOpts

		blockTime   time.Duration
		memPool     *TxPool
		isValidator bool
		rpc         chan RPC
		quit        chan struct{}
	}
)

func NewServer(opts ServerOpts) *Server {
	return &Server{
		ServerOpts:  opts,
		blockTime:   opts.BlockTime,
		memPool:     NewTxPool(),
		isValidator: opts.PrivateKey != nil,
		rpc:         make(chan RPC),
		quit:        make(chan struct{}),
	}
}

func (s *Server) Start() {
	s.initTransports()
	ticker := time.NewTicker(s.blockTime)

free:
	for {
		select {
		case rpc := <-s.rpc:
			fmt.Printf("From: %v \tPayload: %v \n", rpc.From, string(rpc.Payload))
		case <-s.quit:
			break free
		case <-ticker.C:
			if s.isValidator {
				s.createNewBlock()
			}
		}
	}

	fmt.Println("Server shutdown")
}

func (s *Server) handleTransaction(tx *core.Transaction) error {
	if err := tx.Verify(); err != nil {
		return err
	}

	hash := tx.Hash(core.TxHasher{})
	if s.memPool.Has(hash) {
		logrus.WithFields(logrus.Fields{"hash": hash.String()}).Info("transaction already in mempool")
		return nil
	}

	logrus.WithFields(logrus.Fields{"hash": hash.String()}).Info("adding new tx to the mempool")
	return s.memPool.Add(tx)
}

func (s *Server) createNewBlock() error {
	fmt.Println("creating a new block")
	return nil
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

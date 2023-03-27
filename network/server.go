package network

import (
	"fmt"
	"time"

	"github.com/isaqueveras/projectx/core"
	"github.com/isaqueveras/projectx/crypto"
	"github.com/sirupsen/logrus"
)

var defaultBlockTime = 5 * time.Second

type (
	ServerOpts struct {
		RPCHandler RPCHandler
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
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime
	}

	s := &Server{
		blockTime:   opts.BlockTime,
		memPool:     NewTxPool(),
		isValidator: opts.PrivateKey != nil,
		rpc:         make(chan RPC),
		quit:        make(chan struct{}, 1),
	}

	if opts.RPCHandler == nil {
		opts.RPCHandler = NewDefaultRPCHandler(s)
	}

	s.ServerOpts = opts
	return s
}

func (s *Server) Start() {
	s.initTransports()
	ticker := time.NewTicker(s.blockTime)

free:
	for {
		select {
		case rpc := <-s.rpc:
			if err := s.RPCHandler.HandleRPC(rpc); err != nil {
				logrus.Error(err)
			}
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

func (s *Server) ProcessTransaction(from NetAddr, tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{})
	if s.memPool.Has(hash) {
		logrus.WithFields(logrus.Fields{"hash": hash.String()}).Info("transaction already in mempool")
		return nil
	}

	if err := tx.Verify(); err != nil {
		return err
	}

	tx.SetFirstSeen(time.Now().UnixNano())

	logrus.WithFields(logrus.Fields{
		"hash":           hash.String(),
		"mempool_lenght": s.memPool.Len(),
	}).Info("adding new tx to the mempool")

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

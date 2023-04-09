package block_stream

import (
	"context"
	"sync"
	"time"

	"github.com/CronCats/croncat-go/pkg/json_rpc"
	"github.com/sirupsen/logrus"
)

type BlockProvider interface {
	GetLatestBlock() (string, error)
}

type RpcBlockProvider struct {
	Name string
	Url  string
}

func NewRpcBlockProvider(name string, url string) *RpcBlockProvider {
	return &RpcBlockProvider{Name: name, Url: url}
}

func (b *RpcBlockProvider) Run(ctx context.Context, wg *sync.WaitGroup, blocks chan<- interface{}, logger *logrus.Entry) {
	b.fetchBlocks(ctx, wg, blocks, logger)
}

func (b *RpcBlockProvider) getLatestBlock() (interface{}, error) {
	client := &json_rpc.Client{Url: b.Url}

	res, err := client.Call("block", nil, 1)
	if err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (b *RpcBlockProvider) fetchBlocks(ctx context.Context, wg *sync.WaitGroup, blocks chan<- interface{}, logger *logrus.Entry) {
	wg.Add(1)
	go func() {
		ticker := time.NewTicker(5 * time.Second)

		for {
			select {
			case <-ctx.Done():
				logger.WithField("provider", b.Name).Info("Stopping block stream...")
				wg.Done()
				return
			case <-ticker.C:
				logger.WithField("provider", b.Name).Debug("Fetching block")
				block, err := b.getLatestBlock()
				logger.WithField("provider", b.Name).Debug("Got block")
				if err != nil {
					logger.WithField("provider", b.Name).Error("Error fetching block: ", err)
				} else {
					blocks <- block
				}
			}
		}
	}()
}

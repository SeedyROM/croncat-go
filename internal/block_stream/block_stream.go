package block_stream

import (
	"context"
	"errors"
	"sync"

	"github.com/sirupsen/logrus"
)

type BlockStream struct {
	Providers map[string]*RpcBlockProvider

	wg        *sync.WaitGroup
	ctx       context.Context
	logger    *logrus.Entry
	blockChan chan int
	done      chan struct{}
}

func NewBlockStream(ctx context.Context, logger *logrus.Entry) *BlockStream {
	return &BlockStream{
		Providers: make(map[string]*RpcBlockProvider),
		wg:        &sync.WaitGroup{},
		ctx:       ctx,
		logger:    logger,
		blockChan: make(chan int, 128),
		done:      make(chan struct{}, 1),
	}
}

func (b *BlockStream) AddProvider(name string, url string, pollIntervalSecs float64) {
	b.Providers[name] = NewRpcBlockProvider(name, url, pollIntervalSecs)
}

func (b *BlockStream) ProviderStatus(name string) (ProviderStatus, error) {
	provider, ok := b.Providers[name]
	if !ok {
		return ProviderUnknown, errors.New("provider not found")
	}

	return provider.Status, nil
}

func (b *BlockStream) Output() chan int {
	return b.blockChan
}

func (b *BlockStream) Done() <-chan struct{} {
	return b.done
}

func (b *BlockStream) Run() {
	// Run teh providers
	for _, provider := range b.Providers {
		provider.Run(b.ctx, b.wg, b.blockChan, b.logger)
	}

	b.wg.Wait()

	// Close the block channel
	close(b.blockChan)

	// Signal that we are done
	b.done <- struct{}{}
}

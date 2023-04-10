package block_stream

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/CronCats/croncat-go/internal/json_rpc"
	"github.com/sirupsen/logrus"
)

type BlockProvider interface {
	getLatestBlock() (string, error)
}

type ProviderStatus uint64

const (
	ProviderUnknown ProviderStatus = iota
	ProviderOffline
	ProviderIntermittent
	ProviderOnline
)

type RpcBlockProvider struct {
	Name             string
	Url              string
	PollIntervalSecs float64
	Status           ProviderStatus
}

func NewRpcBlockProvider(name string, url string, pollIntervalSecs float64) *RpcBlockProvider {
	return &RpcBlockProvider{
		Name:             name,
		Url:              url,
		PollIntervalSecs: pollIntervalSecs,
		Status:           ProviderUnknown,
	}
}

func (b *RpcBlockProvider) Run(ctx context.Context, wg *sync.WaitGroup, blocks chan<- int, logger *logrus.Entry) {
	b.fetchBlocks(ctx, wg, blocks, logger)
}

func (b *RpcBlockProvider) getLatestBlock() (int, error) {
	client := &json_rpc.Client{Url: b.Url}

	resp, err := client.Call("block", nil, 1)
	if err != nil {
		return 0, err
	}

	return getHeightFromBlockResponse(resp)
}

func (b *RpcBlockProvider) fetchBlocks(ctx context.Context, wg *sync.WaitGroup, blocks chan<- int, logger *logrus.Entry) {
	// Function to get a block from the rpc
	getBlock := func() {
		logger.WithField("provider", b.Name).Debug("Fetching block")
		block, err := b.getLatestBlock()
		logger.WithField("provider", b.Name).Debug("Got block")
		if err != nil {
			b.Status = ProviderOffline
			logger.WithField("provider", b.Name).Debug("Error fetching block: ", err)
		} else {
			b.Status = ProviderOnline
			blocks <- block
		}
	}

	// Block fetch goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()

		logger.WithField("provider", b.Name).Debug("Starting block fetcher")

		ticker := time.NewTicker(time.Duration(b.PollIntervalSecs * float64(time.Second)))

	loop:
		for {
			getBlock()

			select {
			case <-ctx.Done():
				logger.WithField("provider", b.Name).Debug("Stopping block stream...")
				break loop
			case <-ticker.C:
				// Can't seem to get the coverage to work for this case
				//notest
				continue
			}
		}
	}()
}

func getHeightFromBlockResponse(resp *json_rpc.Response) (int, error) {
	result := resp.Result
	block := result.(map[string]interface{})["block"]
	header := block.(map[string]interface{})["header"]
	height, err := strconv.Atoi(header.(map[string]interface{})["height"].(string))

	if err != nil {
		return 0, err
	}

	return height, nil
}

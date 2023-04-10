package block_stream

import (
	"context"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestBlockStream(t *testing.T) {
	defer gock.Off()

	gock.New("http://localhost:1234").
		Post("/").
		Reply(200).
		JSON(map[string]interface{}{
			"jsonrpc": "2.0",
			"result": map[string]interface{}{
				"block": map[string]interface{}{
					"header": map[string]interface{}{
						"height": "1337",
					},
				},
			},
			"id": 1,
		})

	// create a logger instance
	logger := logrus.New().WithField("test", "block_stream")
	logger.Logger.SetOutput(io.Discard)

	// create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10.0*float64(time.Second)))
	defer cancel()

	// create a new block stream
	stream := NewBlockStream(ctx, logger)

	// add a provider
	stream.AddProvider("test_provider", "http://localhost:1234", 0.1)

	// run the block stream
	go stream.Run()

	blockCount := 0

	// test that blocks are received
	select {
	case block := <-stream.Output():
		assert.NotNil(t, block)
		blockCount++
		if blockCount > 1 {
			cancel()
		}
	case <-stream.Done():
		assert.Equal(t, blockCount, 2)
	}
}

func TestNonExistentProvider(t *testing.T) {
	// create a logger instance
	logger := logrus.NewEntry(logrus.StandardLogger())
	logger.Logger.SetOutput(io.Discard)

	// create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(0.1*float64(time.Second)))
	defer cancel()

	// create a new block stream
	stream := NewBlockStream(ctx, logger)

	// check that the provider is marked as offline
	status, err := stream.ProviderStatus("test_provider")
	assert.NotNil(t, err)
	assert.Equal(t, status, ProviderUnknown)
}

func TestBadProvider(t *testing.T) {
	// create a logger instance
	logger := logrus.NewEntry(logrus.StandardLogger())
	logger.Logger.SetOutput(io.Discard)

	// create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(0.01*float64(time.Second)))
	defer cancel()

	// create a new block stream
	stream := NewBlockStream(ctx, logger)

	// add a provider
	stream.AddProvider("test_provider", "http://locaxlhost:1234", 1)

	// run the block stream
	go stream.Run()

	<-ctx.Done()

	// check that the provider is marked as offline
	status, err := stream.ProviderStatus("test_provider")
	assert.Nil(t, err)
	assert.NotNil(t, status)
	assert.Equal(t, status, ProviderOffline)
}

func TestBlockDeDuper(t *testing.T) {
	// create the context
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(0.1*float64(time.Second)))
	defer cancel()

	// create the input channel
	input := make(chan int, 128)

	// create the deduper
	deduper := NewBlockDeDuper(ctx, input)

	// run the deduper
	go deduper.Run()

	// send some values
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= 10; i++ {
			input <- i
			input <- i
			input <- i
		}
		close(input)
	}()

	// verify that the output is correct
	var result []int
	for {
		select {
		case block := <-deduper.Output():
			result = append(result, block)
			if len(result) == 10 {
				cancel()
			}
		case <-ctx.Done():
			wg.Wait()
			assert.Equal(t, result, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
			return
		}
	}
}

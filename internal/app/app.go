package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/CronCats/croncat-go/internal/block_stream"
	"github.com/sirupsen/logrus"
)

type App struct {
	Logger *logrus.Entry
}

func (a *App) Register() {
	a.Logger.Info("Registering croncat tasks")
}

func (a *App) Unregister() {
	a.Logger.Info("Unregistering croncat tasks")
}

func (a *App) Run() {
	a.Logger.Info("Running croncat tasks")

	// Create the application context
	ctx, cancel := context.WithCancel(context.Background())

	// Ctrl+C handler
	go a.ctrlCHandler(cancel)

	// Create the block stream
	stream := block_stream.NewBlockStream(ctx, a.Logger)
	stream.AddProvider("Juno", "https://rpc.uni.junonetwork.io", 3)
	stream.AddProvider("Polkachu", "https://juno-testnet-rpc.polkachu.com", 3)
	stream.AddProvider("Reecepbcups", "https://uni-rpc.reece.sh", 3)

	// De-duplicate the blocks
	deduper := block_stream.NewBlockDeDuper(ctx, stream.Output())

	// Run the block stream pipeline
	go stream.Run()
	go deduper.Run()

	// Local wait group
	wg := &sync.WaitGroup{}
	// Wait for the block stream to finish
	wg.Add(1)
	go func() {
		// Consume the blocks
		for {
			select {
			case <-stream.Done():
				wg.Done()
				return
			case height := <-deduper.Output():

				a.Logger.Infof("Got block %d", height)
			}
		}
	}()
	// Wait for the application to finish
	wg.Wait()

}

func (a *App) ctrlCHandler(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println("")
	a.Logger.Info("SIGINT received, shutting down")
	cancel()
}

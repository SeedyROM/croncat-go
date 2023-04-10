package block_stream

import (
	"context"
	"sync"
)

type BlockDeDuper struct {
	in       chan int
	out      chan int
	lastSeen int

	ctx context.Context
	wg  *sync.WaitGroup
}

func NewBlockDeDuper(ctx context.Context, in chan int) *BlockDeDuper {
	return &BlockDeDuper{
		in:  in,
		out: make(chan int, 128),
		ctx: ctx,
		wg:  &sync.WaitGroup{},
	}
}

func (b *BlockDeDuper) Run() {
	b.wg.Add(1)
	defer b.wg.Done()

	for {
		select {
		case block := <-b.in:
			if block > b.lastSeen {
				b.out <- block
				b.lastSeen = block
			}
		case <-b.ctx.Done():
			return
		}
	}
}

func (b *BlockDeDuper) Output() chan int {
	return b.out
}

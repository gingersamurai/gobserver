package closer

import (
	"context"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Closer struct {
	sync.Mutex
	closeTimeout time.Duration
	funcs        []func(ctx context.Context) error
}

func NewCloser(closeTimeout time.Duration) *Closer {
	return &Closer{closeTimeout: closeTimeout}
}

func (c *Closer) Add(f func(ctx context.Context) error) {
	c.Lock()
	defer c.Unlock()

	c.funcs = append(c.funcs, f)
}

func (c *Closer) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), c.closeTimeout)
	defer cancel()
	doneCh := make(chan struct{})
	go func() {
		var wg sync.WaitGroup
		wg.Add(len(c.funcs))
		for _, f := range c.funcs {
			f := f
			go func() {
				defer wg.Done()

				if err := f(ctx); err != nil {
					log.Println(err)
				}
			}()
		}

		wg.Wait()
		doneCh <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		log.Fatal(ctx.Err())
	case <-doneCh:
		break
	}

}

func (c *Closer) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	log.Println("got exit signal, starting graceful shutdown")
	c.Close()
	log.Println("finished graceful shutdown")

}

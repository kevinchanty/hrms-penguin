package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/gen2brain/beeep"
)

const defaultTick = 15 * time.Second

type config struct {
	tick time.Duration
}

func (c *config) init(args []string) error {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)

	tick := flags.Duration("tick", defaultTick, "Ticking interval")

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	c.tick = *tick
	return nil
}

func main() {
	fmt.Println("Start ticking")

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	c := &config{}

	defer func() {
		cancel()
	}()

	if err := run(ctx, c); err != nil {
		os.Exit(1)
	}
}

func run(ctx context.Context, c *config) error {
	c.init(os.Args)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.Tick(c.tick):
			beeep.Notify("Kirby is talking", "Poyo!", "")
		}
	}
}

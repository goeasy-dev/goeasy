package main

import (
	"context"
	"fmt"
	"time"

	"goeasy.dev"
	"goeasy.dev/bootstrap"
	"goeasy.dev/status"
	"goeasy.dev/status/statustype"
)

func main() {
	ctx, done := bootstrap.Bootstrap()
	defer done()

	application := goeasy.NewApplication()
	statusCheck := status.SimpleCheck(statustype.Startup & statustype.Readiness)

	toggleRunner := func(ctx context.Context) (goeasy.StopFunc, error) {
		done := make(chan struct{})

		t := time.NewTicker(time.Second * 5)
		go func() {
			for {
				select {
				case <-done:
					t.Stop()
					return
				case <-t.C:
					*statusCheck = !*statusCheck
				}
			}
		}()

		return func(ctx context.Context) error {
			close(done)
			return nil
		}, nil
	}

	err := application.Start(ctx, toggleRunner)
	if err != nil {
		fmt.Println(err)
	}
}

package runners

import (
	"context"
	"log"
	"time"

	"goeasy.dev"
	"goeasy.dev/errors"
)

type TimedOptions struct {
	HandleFunc func(ctx context.Context) error
	CancelFunc func(ctx context.Context) error
	Interval   time.Duration
}

// Timed executes `TimedOptions.HandleFunc` at a given interval specified by `TimedOptions.Interval`
// If `TimedOptions.HandleFunc` returns an error, the error is logged and the timer resets.
func Timed(options TimedOptions) goeasy.RunnerFunc {
	return func(ctx context.Context) (goeasy.StopFunc, error) {
		done := make(chan struct{})

		go func() {
			t := time.NewTimer(options.Interval)
			for {
				select {
				case <-done:
					return
				case <-t.C:
				}

				err := options.HandleFunc(context.Background())
				if err != nil {
					log.Println(errors.Wrap(err, "error executing timed runner function"))
				}

				t.Reset(options.Interval)
			}
		}()

		return func(ctx context.Context) error {
			close(done)
			var err error
			if options.CancelFunc != nil {
				err = options.CancelFunc(ctx)
			}

			return err
		}, nil
	}
}

package goeasy

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"goeasy.dev/util"
)

type StopFunc func(ctx context.Context) error

type RunnerFunc func(ctx context.Context) (StopFunc, error)

type Application struct {
	runners []RunnerFunc
}

var defaultRunners = make([]RunnerFunc, 0, 10)

func RegisterDefaultRunner(runner RunnerFunc) {
	defaultRunners = append(defaultRunners, runner)
}

func NewApplication() *Application {
	return &Application{
		runners: append([]RunnerFunc{}, defaultRunners...),
	}
}

func (a *Application) RegisterRunnerFunc(runner RunnerFunc) {
	a.runners = append(a.runners, runner)
}

func (a Application) Start(ctx context.Context, runners ...RunnerFunc) error {
	runners = append(a.runners, runners...)
	stopFuncs := make([]StopFunc, len(runners))

	var err error
	for i, runner := range runners {
		stopFuncs[i], err = runner(ctx)
		if err != nil {
			return fmt.Errorf("unable to start runner: %w", err)
		}
	}

	log.Println("all runners started")

	waitForStopSignal(ctx)

	// TODO: implement timeout for application stop
	// TODO: implement timeout for runner stop
	for _, stopFunc := range util.ReverseSlice(stopFuncs) {
		err = stopFunc(ctx)
		if err != nil {
			log.Println("unable to stop runner: %w", err)
		}
	}

	log.Println("all runners stopped")

	return nil
}

func waitForStopSignal(ctx context.Context) {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT)

	select {
	case <-sigs:
	case <-ctx.Done():
	}

	log.Println("stop signal received")
}

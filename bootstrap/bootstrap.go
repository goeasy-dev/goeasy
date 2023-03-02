package bootstrap

import (
	"context"

	"goeasy.dev"
	"goeasy.dev/internal/control"
)

func Bootstrap() (context.Context, func()) {
	ctx, cancel := context.WithCancel(context.Background())

	goeasy.RegisterDefaultRunner(control.NewRunner())

	return ctx, cancel
}

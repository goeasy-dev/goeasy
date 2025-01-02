package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"goeasy.dev/cmd/goeasy/commands/gen"
)

func NewGoeasyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "goeasy",
		Short: "GoEasy is a framework for building applications in Go",
	}

	cmd.AddCommand(gen.NewGenCommand())

	return cmd
}

func Execute() {
	if err := NewGoeasyCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

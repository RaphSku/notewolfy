package cmd

import (
	"context"
	"os"

	"github.com/RaphSku/notewolfy/cmd/version"
	"github.com/RaphSku/notewolfy/internal/console"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type CLI struct {
	ctx    context.Context
	logger *zap.Logger

	rootCmd *cobra.Command
}

func NewCLI(ctx context.Context, logger *zap.Logger) *CLI {
	rootCmd := &cobra.Command{
		Use:   "notewolfy",
		Short: "notewolfy a minimalistic note taking console application",
		Long:  `notewolfy is a minimalistic note taking console application that allows you to organize and manage your markdown notes`,
		Run: func(cmd *cobra.Command, args []string) {
			console.StartConsoleApplication(ctx)
		},
	}

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	return &CLI{
		ctx:     ctx,
		logger:  logger,
		rootCmd: rootCmd,
	}
}

func (cli *CLI) AddSubCommands() {
	versionCmd := version.NewVersionCommand()
	cli.rootCmd.AddCommand(versionCmd)
}

func (cli *CLI) Execute() {
	if err := cli.rootCmd.Execute(); err != nil {
		cli.logger.Info("CLI failed to run", zap.Error(err))
		os.Exit(1)
	}
}

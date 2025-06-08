package cmd

import (
	"fmt"
	"os"

	"github.com/RaphSku/notewolfy/cmd/version"
	"github.com/RaphSku/notewolfy/internal/console"
	"github.com/RaphSku/notewolfy/internal/logging"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type CLI struct {
	logger *zap.Logger

	rootCmd *cobra.Command
}

func NewCLI() *CLI {
	return &CLI{}
}

func (cli *CLI) Run() {
	rootCmd := &cobra.Command{
		Use:   "notewolfy",
		Short: "notewolfy, a minimalistic note taking console application!",
		Long:  `notewolfy is a minimalistic note taking console application that allows you to organize and manage your markdown notes!`,
		Run:   cli.runNotewolfyCommand,
	}

	// --- ROOT CMD
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	cli.rootCmd = rootCmd

	// --- SUB CMD
	versionCmd := version.NewVersionCmd().GetVersionCmd()
	cli.rootCmd.AddCommand(versionCmd)

	// --- EXECUTE
	if err := cli.rootCmd.Execute(); err != nil {
		fmt.Printf("CLI failed to run due to %v\n", err)
	}
}

func (cli *CLI) runNotewolfyCommand(cmd *cobra.Command, args []string) {
	// if debugLevel = 1, then debug logs are shown to os.Stdout, otherwise no logs will be printed
	debugLevel := os.Getenv("NOTEWOLFY_DEBUG")
	var logger *zap.Logger
	if debugLevel == "1" {
		logger = logging.SetupZapLogger(true)
	} else {
		logger = logging.SetupZapLogger(false)
	}
	cli.logger = logger

	console.StartConsoleApplication()
}

package main

import (
	"context"
	"os"

	"github.com/RaphSku/notewolfy/cmd"
	"github.com/RaphSku/notewolfy/internal/logging"
	"go.uber.org/zap"
)

func main() {
	// if debugLevel = 1, then debug logs are shown to os.Stdout, otherwise no logs will be printed
	debugLevel := os.Getenv("NOTEWOLFY_DEBUG")
	var logger *zap.Logger
	if debugLevel == "1" {
		logger = logging.SetupZapLogger(true)
	} else {
		logger = logging.SetupZapLogger(false)
	}

	cli := cmd.NewCLI(context.Background(), logger)
	cli.AddSubCommands()
	cli.Execute()
}

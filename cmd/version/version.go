package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

const VERSION = "v0.1.0"

func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version of notewolfy.",
		Long:  `This will show you the version of notewolfy in the format: {MAJOR}-{MINOR}-{PATCH}.`,
		Run:   versionCommandFunc,
	}
}

func versionCommandFunc(cmd *cobra.Command, args []string) {
	fmt.Println(VERSION)
}

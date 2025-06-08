package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

const VERSION = "v0.2.0"

type VersionCmd struct{}

func NewVersionCmd() *VersionCmd {
	return &VersionCmd{}
}

func (vc *VersionCmd) GetVersionCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Prints the version number of notewolfy.",
		Long:  `This will show you the version of notewolfy in the format: v{MAJOR}.{MINOR}.{PATCH}.`,
		Run:   vc.runVersionCmd,
	}

	return versionCmd
}

func (vc *VersionCmd) runVersionCmd(cmd *cobra.Command, args []string) {
	fmt.Println(VERSION)
}

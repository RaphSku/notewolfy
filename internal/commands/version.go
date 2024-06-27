package commands

import (
	"fmt"

	"github.com/RaphSku/notewolfy/cmd/version"
)

type VersionStrategy struct{}

func (vs *VersionStrategy) Run() error {
	fmt.Printf("\n\rnotewolfy version %s at your disposal!", version.VERSION)
	return nil
}

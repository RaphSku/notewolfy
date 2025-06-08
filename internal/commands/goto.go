package commands

import (
	"fmt"
	"regexp"

	"github.com/RaphSku/notewolfy/internal/structure"
)

type GoToStrategy struct {
	statement string
	mmf       *structure.MetadataNoteWolfyFileHandle
}

func (gts *GoToStrategy) Run() error {
	nameCaptureGroupName := "name"
	nodeNamePattern := "[\\w]+"
	pattern := fmt.Sprintf("goto (?P<%s>%s)", nameCaptureGroupName, nodeNamePattern)
	goToRegex := regexp.MustCompile(pattern)
	matches := goToRegex.FindStringSubmatch(gts.statement)
	if len(matches) != 2 {
		return fmt.Errorf("\n\rPlease check whether the node matches the regex %s!", nodeNamePattern)
	}
	names := goToRegex.SubexpNames()
	var goToName string
	for i, name := range names[1:] {
		if name == nameCaptureGroupName {
			goToName = matches[i+1]
		}
	}

	activeNodeName := gts.mmf.ActiveNode
	activeNode := gts.mmf.FindNode(activeNodeName)
	for _, childNode := range activeNode.Children {
		if childNode.Name == goToName {
			gts.mmf.ActiveNode = childNode.Name
			gts.mmf.Save()
			return nil
		}
	}

	return fmt.Errorf("\r\nCould not find node '%s'!", goToName)
}

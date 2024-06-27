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
	goToRegex := regexp.MustCompile("goto (?P<name>[[:alpha:]]+)")
	matches := goToRegex.FindStringSubmatch(gts.statement)
	names := goToRegex.SubexpNames()
	namedGroups := make(map[string]string)
	for i, name := range names {
		if i != 0 && name != "" {
			namedGroups[name] = matches[i]
		}
	}
	goToName := namedGroups["name"]

	activeNodeName := gts.mmf.ActiveNode
	activeNode := gts.mmf.FindNode(activeNodeName)
	for _, childNode := range activeNode.Children {
		if childNode.Name == goToName {
			gts.mmf.ActiveNode = childNode.Name
			gts.mmf.Save()
			return nil
		}
	}

	fmt.Println("\r\nCould not find node ", goToName)

	return nil
}

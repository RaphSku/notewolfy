package commands

import (
	"regexp"

	"github.com/RaphSku/notewolfy/internal/structure"
)

type OpenStrategy struct {
	statement string
	mmf       *structure.MetadataNoteWolfyFileHandle
}

func (ops *OpenStrategy) Run() error {
	workspaceNameRegex := regexp.MustCompile("open (?P<name>[[:alpha:]]+)")
	matches := workspaceNameRegex.FindStringSubmatch(ops.statement)
	names := workspaceNameRegex.SubexpNames()
	namedGroups := make(map[string]string)
	for i, name := range names {
		if i != 0 && name != "" {
			namedGroups[name] = matches[i]
		}
	}
	workspaceName := namedGroups["name"]

	for _, workspace := range ops.mmf.Workspaces {
		if workspace.Name == workspaceName {
			ops.mmf.ActiveWorkspace = workspace.Name
			ops.mmf.ActiveNode = workspace.Name
			ops.mmf.Save()
		}
	}

	return nil
}

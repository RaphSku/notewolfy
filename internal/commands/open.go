package commands

import (
	"fmt"
	"regexp"

	"github.com/RaphSku/notewolfy/internal/structure"
)

type OpenStrategy struct {
	statement string
	mmf       *structure.MetadataNoteWolfyFileHandle
}

func (ops *OpenStrategy) Run() error {
	nameCaptureGroupName := "name"
	workspaceNamePattern := "[\\w]+"
	pattern := fmt.Sprintf("open (?P<%s>%s)", nameCaptureGroupName, workspaceNamePattern)
	workspaceNameRegex := regexp.MustCompile(pattern)
	matches := workspaceNameRegex.FindStringSubmatch(ops.statement)
	if len(matches) != 2 {
		return fmt.Errorf("\n\rPlease check whether the workspace name matches the regex %s!", workspaceNamePattern)
	}
	names := workspaceNameRegex.SubexpNames()
	var workspaceName string
	for i, name := range names[1:] {
		if name == nameCaptureGroupName {
			workspaceName = matches[i+1]
		}
	}

	foundWorkspace := false
	for _, workspace := range ops.mmf.Workspaces {
		if workspace.Name == workspaceName {
			ops.mmf.ActiveWorkspace = workspace.Name
			ops.mmf.ActiveNode = workspace.Name
			ops.mmf.Save()

			foundWorkspace = true
		}
	}

	if !foundWorkspace {
		return fmt.Errorf("\n\rDid not find workspace '%s'! Please specify an existing workspace.", workspaceName)
	}

	return nil
}

package commands

import (
	"fmt"
	"os"
	"regexp"

	"github.com/RaphSku/notewolfy/internal/structure"
	"github.com/RaphSku/notewolfy/internal/utility"
)

type CreateWorkspaceStrategy struct {
	statement string
	mmf       *structure.MetadataNoteWolfyFileHandle
}

func (cws *CreateWorkspaceStrategy) Run() error {
	pathRegex := regexp.MustCompile("[~.]{0,1}/.*")
	workspacePath := pathRegex.FindAllString(cws.statement, 1)
	if len(workspacePath) == 0 {
		return nil
	}

	workspaceNameRegex := regexp.MustCompile("create workspace (?P<name>[\\w]+) [~/.]{0,1}/.*")
	matches := workspaceNameRegex.FindStringSubmatch(cws.statement)
	names := workspaceNameRegex.SubexpNames()
	namedGroups := make(map[string]string)
	for i, name := range names {
		if i != 0 && name != "" {
			namedGroups[name] = matches[i]
		}
	}
	workspaceName := namedGroups["name"]

	pathToWorkspace, err := utility.ExpandRelativePaths(workspacePath[0])
	if err != nil {
		return err
	}
	cws.mmf.AddNewWorkspace(workspaceName, pathToWorkspace)
	cws.mmf.Save()

	err = os.Mkdir(pathToWorkspace, 0755)
	if err != nil {
		return err
	}

	return nil
}

type DeleteWorkspaceStrategy struct {
	statement string
	mmf       *structure.MetadataNoteWolfyFileHandle
}

func (dws *DeleteWorkspaceStrategy) Run() error {
	workspaceNameRegex := regexp.MustCompile("delete workspace (?P<name>[\\w]+)")
	matches := workspaceNameRegex.FindStringSubmatch(dws.statement)
	names := workspaceNameRegex.SubexpNames()
	namedGroups := make(map[string]string)
	for i, name := range names {
		if i != 0 && name != "" {
			namedGroups[name] = matches[i]
		}
	}
	workspaceName := namedGroups["name"]

	foundIndex := -1
	for index, workspace := range dws.mmf.Workspaces {
		if workspace.Name == workspaceName {
			foundIndex = index
		}
	}
	if foundIndex == -1 {
		return fmt.Errorf("workspace %s could not be found!", workspaceName)
	}

	if len(dws.mmf.Workspaces[foundIndex].Children) != 0 || len(dws.mmf.Workspaces[foundIndex].Markdowns) != 0 {
		return fmt.Errorf("before you delete a workspace, ensure that you have deleted all nodes and markdown files in this workspace!")
	}

	workspacePath := dws.mmf.Workspaces[foundIndex].Path

	dws.mmf.Workspaces = append(dws.mmf.Workspaces[:foundIndex], dws.mmf.Workspaces[foundIndex+1:]...)
	if len(dws.mmf.Workspaces) != 0 {
		dws.mmf.ActiveWorkspace = dws.mmf.Workspaces[0].Name
		dws.mmf.ActiveNode = dws.mmf.Workspaces[0].Name
	} else {
		dws.mmf.ActiveWorkspace = ""
		dws.mmf.ActiveNode = ""
	}
	dws.mmf.Save()

	err := os.Remove(workspacePath)
	if err != nil {
		return fmt.Errorf("the workspace %s could not be deleted, please clean up the following workspace path yourself: %s, error: %v", workspaceName, workspacePath, err)
	}

	return nil
}

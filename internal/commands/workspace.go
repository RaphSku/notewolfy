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
	nameCaptureGroupName := "name"
	pathCaptureGroupName := "path"
	workspaceNamePattern := "[\\w]+"
	workspacePathPattern := "[~.]{0,1}/{0,1}.*"
	pattern := fmt.Sprintf("create workspace (?P<%s>%s) (?P<%s>%s)", nameCaptureGroupName, workspaceNamePattern, pathCaptureGroupName, workspacePathPattern)
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(cws.statement)
	if len(matches) != 3 {
		return fmt.Errorf("\n\rPlease check whether the workspace name matches the regex %s and whether the workspace path matches the regex %s!", workspaceNamePattern, workspacePathPattern)
	}
	names := regex.SubexpNames()
	var workspaceName string
	var workspacePath string
	for i, name := range names[1:] {
		if name == nameCaptureGroupName {
			workspaceName = matches[i+1]
		} else if name == pathCaptureGroupName {
			workspacePath = matches[i+1]
		}
	}

	pathToWorkspace, err := utility.ExpandRelativePaths(workspacePath)
	if err != nil {
		return err
	}

	fileInfo, err := os.Stat(pathToWorkspace)
	if err == nil {
		if fileInfo.IsDir() {
			return fmt.Errorf("\n\rThe specified path %s already exists, please choose a valid path!", pathToWorkspace)
		}
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
	nameCaptureGroupName := "name"
	workspaceNamePattern := "[\\w]+"
	pattern := fmt.Sprintf("delete workspace (?P<%s>%s)", nameCaptureGroupName, workspaceNamePattern)
	workspaceNameRegex := regexp.MustCompile(pattern)
	matches := workspaceNameRegex.FindStringSubmatch(dws.statement)
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

	foundIndex := -1
	for index, workspace := range dws.mmf.Workspaces {
		if workspace.Name == workspaceName {
			foundIndex = index
		}
	}
	if foundIndex == -1 {
		return fmt.Errorf("\n\rWorkspace '%s' could not be found!", workspaceName)
	}

	if len(dws.mmf.Workspaces[foundIndex].Children) != 0 || len(dws.mmf.Workspaces[foundIndex].Markdowns) != 0 {
		return fmt.Errorf("\n\rBefore you delete a workspace, ensure that you have deleted all nodes and markdown files in this workspace!")
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
		return fmt.Errorf("\n\rThe workspace '%s' could not be deleted, please clean up the following workspace path yourself: %s, error: %v", workspaceName, workspacePath, err)
	}
	fmt.Printf("\n\rDeleted workspace '%s' successfully!", workspaceName)

	return nil
}

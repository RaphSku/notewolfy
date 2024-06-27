package commands

import (
	"fmt"

	"github.com/RaphSku/notewolfy/internal/structure"
)

type ListStrategy struct {
	mmf *structure.MetadataNoteWolfyFileHandle
}

func (ls *ListStrategy) Run() error {
	activeNodeName := ls.mmf.ActiveNode
	if activeNodeName == "" {
		fmt.Print("\n\rSeems like you have not created a workspace yet! Create one with 'create workspace <workspace_name> <workspace_path>'")
		return nil
	}
	activeNode := ls.mmf.FindNode(activeNodeName)
	if activeNode == nil {
		fmt.Print("\n\rSeems like you have not created a workspace yet! At least no active node is set!")
	}
	ls.mmf.ListResourcesOnNode(activeNode)

	return nil
}

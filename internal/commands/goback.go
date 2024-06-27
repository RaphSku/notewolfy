package commands

import (
	"github.com/RaphSku/notewolfy/internal/structure"
)

type GoBackStrategy struct {
	mmf *structure.MetadataNoteWolfyFileHandle
}

func (gbs *GoBackStrategy) Run() error {
	activeNodeName := gbs.mmf.ActiveNode
	parentNode := gbs.mmf.FindParentNode(activeNodeName)
	if parentNode != nil {
		gbs.mmf.ActiveNode = parentNode.Name
		gbs.mmf.Save()
	}

	return nil
}

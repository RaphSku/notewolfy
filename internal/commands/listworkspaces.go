package commands

import (
	"github.com/RaphSku/notewolfy/internal/structure"
)

type ListWorkspacesStrategy struct {
	statement string
	mmf       *structure.MetadataNoteWolfyFileHandle
}

func (lws *ListWorkspacesStrategy) Run() error {
	lws.mmf.ListWorkspaces()

	return nil
}

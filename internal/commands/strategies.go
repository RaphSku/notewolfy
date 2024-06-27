package commands

import (
	"strings"

	"github.com/RaphSku/notewolfy/internal/structure"
)

func matchStatementToStrategy(mmf *structure.MetadataNoteWolfyFileHandle, statement string) Strategy {
	strategies := map[string]Strategy{
		"version": &VersionStrategy{},
		"create workspace": &CreateWorkspaceStrategy{
			statement: statement,
			mmf:       mmf,
		},
		"delete workspace": &DeleteWorkspaceStrategy{
			statement: statement,
			mmf:       mmf,
		},
		"ls ws": &ListWorkspacesStrategy{
			statement: statement,
			mmf:       mmf,
		},
		"create node": &CreateNodeStrategy{
			statement: statement,
			mmf:       mmf,
		},
		"delete node": &DeleteNodeStrategy{
			statement: statement,
			mmf:       mmf,
		},
		"create md": &CreateMarkdownStrategy{
			statement: statement,
			mmf:       mmf,
		},
		"delete md": &DeleteMDStrategy{
			statement: statement,
			mmf:       mmf,
		},
		"open": &OpenStrategy{
			statement: statement,
			mmf:       mmf,
		},
		"edit": &EditStrategy{
			statement: statement,
			mmf:       mmf,
		},
		"ls": &ListStrategy{
			mmf: mmf,
		},
		"goto": &GoToStrategy{
			statement: statement,
			mmf:       mmf,
		},
		"goback": &GoBackStrategy{
			mmf: mmf,
		},
		"help": &HelpStrategy{
			statement: statement,
		},
	}

	var longestPrefixMatch string
	longestPrefixLength := 0
	for command := range strategies {
		if len(command) <= len(statement) {
			matches := true
			for i := range command {
				if statement[i] != command[i] {
					matches = false
					break
				}
			}
			if matches && len(command) > longestPrefixLength {
				if len(command) == len(statement) || (len(command) < len(statement) && strings.HasPrefix(string(statement[len(command)]), " ")) {
					longestPrefixMatch = command
					longestPrefixLength = len(command)
				}
			}
		}
	}

	if longestPrefixLength == 0 {
		return nil
	}

	return strategies[longestPrefixMatch]
}

package commands

import (
	"fmt"
	"regexp"
)

type HelpStrategy struct {
	statement string
}

func (hs *HelpStrategy) Run() error {
	helpRegex := regexp.MustCompile("^help (?P<name>[[:alpha:]]+(?: [[:alpha:]]+)*)")
	matches := helpRegex.FindStringSubmatch(hs.statement)
	names := helpRegex.SubexpNames()
	namedGroups := make(map[string]string)
	for i, name := range names {
		if i != 0 && name != "" {
			namedGroups[name] = matches[i]
		}
	}
	helpCommand := namedGroups["name"]

	var command string
	var description string
	var example string
	switch helpCommand {
	case "ls":
		command = "\n\rCommand: ls"
		description = "\n\rDescription: ls can be used to list information about the node that you are on, e.g. active node, markdown files on that node, etc."
		example = "\n\rExample Usage: ls"
	case "ls ws":
		command = "\n\rCommand: ls ws"
		description = "\n\rDescription: ls ws will list the workspaces and their root paths in a table format."
		example = "\n\rExample Usage: ls ws"
	case "create workspace":
		command = "\n\rCommand: create workspace <workspaceName> <workspacePath>"
		description = "\n\rDescription: create workspace will create a new workspace for you under the specified name and path that you can choose."
		example = "\n\rExample Usage: create workspace example /path/to/example"
	case "delete workspace":
		command = "\n\rCommand: delete workspace <workspaceName>"
		description = "\n\rDescription: delete workspace lets you delete the specified workspace. This will fail if nodes & markdown files still exist on the node."
		example = "\n\rExample Usage: delete workspace example"
	case "create node":
		command = "\n\rCommand: create node <nodeName>"
		description = "\n\rDescription: create node will create a new node for you under the specified name. The node path will correspond to /pathOfActiveNode/nodeName."
		example = "\n\rExample Usage: create node example"
	case "delete node":
		command = "\n\rCommand: delete node <nodeName>"
		description = "\n\rDescription: delete node lets you delete the specified node. This will fail if markdown files still exist on the node."
		example = "\n\rExample Usage: delete node example"
	case "create md":
		command = "\n\rCommand: create md <markdownFileName>"
		description = "\n\rDescription: create md will create a new markdown file for you under the specified name. You don't need to append the file extension to the name."
		example = "\n\rExample Usage: create md example"
	case "delete md":
		command = "\n\rCommand: delete md <markdownFileName>"
		description = "\n\rDescription: delete md lets you delete the specified markdown file. Specify only the name, so without the file extension."
		example = "\n\rExample Usage: delete md example"
	case "goto":
		command = "\n\rCommand: goto <nodeName>"
		description = "\n\rDescription: goto will let you change the node, specify the name of the node that is a direct child of the node that you are on."
		example = "\n\rExample Usage: goto example"
	case "goback":
		command = "\n\rCommand: goback"
		description = "\n\rDescription: goback lets you go to the parent node of the node that you are currently on."
		example = "\n\rExample Usage: goback"
	case "open":
		command = "\n\rCommand: open <workspaceName>"
		description = "\n\rDescription: open lets you open another workspace, in the sense that the active node will be set to the specified workspace node."
		example = "\n\rExample Usage: open example"
	case "version":
		command = "\n\rCommand: version"
		description = "\n\rDescription: version will print notewolfy's version."
		example = "\n\rExample Usage: version"
	}
	fmt.Printf(command + description + example)

	return nil
}

package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/RaphSku/notewolfy/internal/structure"
	"github.com/RaphSku/notewolfy/internal/utility"
)

type CreateNodeStrategy struct {
	statement string
	mmf       *structure.MetadataNoteWolfyFileHandle
}

func (cns *CreateNodeStrategy) Run() error {
	nodeNameRegex := regexp.MustCompile("create node (?P<name>[\\w]+)")
	matches := nodeNameRegex.FindStringSubmatch(cns.statement)
	names := nodeNameRegex.SubexpNames()
	namedGroups := make(map[string]string)
	for i, name := range names {
		if i != 0 && name != "" {
			namedGroups[name] = matches[i]
		}
	}
	nodeName := namedGroups["name"]

	activeNodeName := cns.mmf.ActiveNode
	activeNode := cns.mmf.FindNode(activeNodeName)
	pathToNode, err := utility.ExpandRelativePaths(filepath.Join(activeNode.Path, nodeName))
	if err != nil {
		return err
	}

	var children []*structure.Node
	var markdowns []*structure.Markdown
	childNode := &structure.Node{
		Name:      nodeName,
		Path:      pathToNode,
		Markdowns: markdowns,
		Children:  children,
	}
	err = cns.mmf.AddChild(childNode)
	if err != nil {
		return err
	}
	cns.mmf.Save()

	err = os.Mkdir(pathToNode, 0750)
	if err != nil {
		return err
	}

	return nil
}

type DeleteNodeStrategy struct {
	statement string
	mmf       *structure.MetadataNoteWolfyFileHandle
}

func (dns *DeleteNodeStrategy) Run() error {
	nodeNameRegex := regexp.MustCompile("delete node (?P<name>[\\w]+)")
	matches := nodeNameRegex.FindStringSubmatch(dns.statement)
	names := nodeNameRegex.SubexpNames()
	namedGroups := make(map[string]string)
	for i, name := range names {
		if i != 0 && name != "" {
			namedGroups[name] = matches[i]
		}
	}
	nodeName := namedGroups["name"]

	activeNodeName := dns.mmf.ActiveNode
	activeNode := dns.mmf.FindNode(activeNodeName)
	for index, child := range activeNode.Children {
		if child.Name == nodeName {
			if len(child.Markdowns) != 0 || len(child.Children) != 0 {
				return fmt.Errorf("please delete all subsequent nodes and markdown files before deleting %s", nodeName)
			}

			err := dns.mmf.DeleteChildByIndex(index)
			if err != nil {
				return err
			}
			err = dns.mmf.Save()
			if err != nil {
				return err
			}

			err = os.Remove(child.Path)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return fmt.Errorf("there is no node with the name %s", nodeName)
}

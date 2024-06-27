package structure

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/RaphSku/notewolfy/internal/utility"
)

type Config struct {
	MetadataFilePath string
}

type Markdown struct {
	Filename string `json:"filename"`
}

type Node struct {
	Name      string      `json:"name"`
	Path      string      `json:"path"`
	Markdowns []*Markdown `json:"markdowns"`
	Children  []*Node     `json:"children"`
}

type MetadataNoteWolfyFileHandle struct {
	Config          *Config `json:"-"`
	Workspaces      []*Node `json:"workspaces"`
	ActiveWorkspace string  `json:"activeworkspace"`
	ActiveNode      string  `json:"activenode"`
}

func NewMetadataNoteWolfyFileHandle(config *Config) (*MetadataNoteWolfyFileHandle, error) {
	notewolfyFileHandle := &MetadataNoteWolfyFileHandle{Config: config}
	err := notewolfyFileHandle.load()
	if err != nil {
		return nil, err
	}
	return notewolfyFileHandle, nil
}

func (mmf *MetadataNoteWolfyFileHandle) AddNewWorkspace(workspaceName string, workspacePath string) error {
	expanedWorkspacePath, err := utility.ExpandRelativePaths(workspacePath)
	if err != nil {
		return err
	}
	var nodes []*Node
	var markdowns []*Markdown
	newWorkspaceNode := &Node{
		Name:      workspaceName,
		Path:      expanedWorkspacePath,
		Markdowns: markdowns,
		Children:  nodes,
	}
	mmf.Workspaces = append(mmf.Workspaces, newWorkspaceNode)
	mmf.ActiveWorkspace = workspaceName
	mmf.ActiveNode = workspaceName

	return nil
}

func (mmf *MetadataNoteWolfyFileHandle) DoesWorkspaceExist(name string) bool {
	doesExist := false
	for _, node := range mmf.Workspaces {
		if node.Name == name {
			doesExist = true
		}
	}

	return doesExist
}

func (mmf *MetadataNoteWolfyFileHandle) Save() error {
	file, err := mmf.getMetadataFile()
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(mmf); err != nil {
		return err
	}
	return nil
}

func (mmf *MetadataNoteWolfyFileHandle) AddChild(childNode *Node) error {
	activeNodeName := mmf.ActiveNode
	if activeNodeName == "" {
		return errors.New("no active node, seems like you have not created a workspace yet")
	}

	activeNode := mmf.FindNode(activeNodeName)

	activePath := activeNode.Path
	childPath := childNode.Path
	_, err := utility.DoesChildPathMatchesParentPath(activePath, childPath)
	if err != nil {
		return err
	}
	activeNode.Children = append(activeNode.Children, childNode)

	return nil
}

func (mmf *MetadataNoteWolfyFileHandle) DeleteChildByIndex(index int) error {
	activeNodeName := mmf.ActiveNode
	if activeNodeName == "" {
		return errors.New("no active node, seems like you have not created a workspace yet")
	}

	activeNode := mmf.FindNode(activeNodeName)
	if index < 0 || index >= len(activeNode.Children) {
		return errors.New("index is out of range, check that the child at this index exists")
	}
	activeNode.Children = append(activeNode.Children[:index], activeNode.Children[index+1:]...)

	return nil
}

func (mmf *MetadataNoteWolfyFileHandle) AddMarkdown(markdown *Markdown) error {
	activeNodeName := mmf.ActiveNode
	if activeNodeName == "" {
		return errors.New("no active node, seems like you have not created a workspace yet")
	}

	activeNode := mmf.FindNode(activeNodeName)
	activeNode.Markdowns = append(activeNode.Markdowns, markdown)

	return nil
}

func (mmf *MetadataNoteWolfyFileHandle) DeleteMarkdown(markdownName string) error {
	activeNodeName := mmf.ActiveNode
	if activeNodeName == "" {
		return errors.New("no active node, seems like you have not created a workspace yet")
	}

	activeNode := mmf.FindNode(activeNodeName)

	foundIndex := -1
	for index, markdown := range activeNode.Markdowns {
		if markdown.Filename[:len(markdown.Filename)-3] == markdownName {
			foundIndex = index
			break
		}
	}

	if foundIndex == -1 {
		return fmt.Errorf("markdownName %s matches no name of a markdown note", markdownName)
	}

	activeNode.Markdowns = append(activeNode.Markdowns[:foundIndex], activeNode.Markdowns[foundIndex+1:]...)
	return nil
}

func (mmf *MetadataNoteWolfyFileHandle) ListWorkspaces() {
	longestStringLength := len("Workspace Name")
	for _, workspace := range mmf.Workspaces {
		workspaceNameLength := len(workspace.Name)
		if workspaceNameLength > longestStringLength {
			longestStringLength = workspaceNameLength
		}
		workspacePathLength := len(workspace.Path)
		if workspacePathLength > longestStringLength {
			longestStringLength = workspacePathLength
		}
	}

	fmt.Printf(fmt.Sprintf("\r\n%%-%[1]d.%[1]ds%%-%[1]d.%[1]ds", longestStringLength), "Workspace Name", "Workspace Path")
	totalWidth := 2*longestStringLength + 1
	fmt.Printf("\r\n%s\n", strings.Repeat("-", totalWidth))
	for _, workspace := range mmf.Workspaces {
		fmt.Printf(fmt.Sprintf("\r%%-%[1]d.%[1]ds%%-%[1]d.%[1]ds\n", longestStringLength), workspace.Name, workspace.Path)
	}
}

func (mmf *MetadataNoteWolfyFileHandle) ListResourcesOnNode(node *Node) {
	fmt.Println("\r\nYou are on node: ", node.Name)
	fmt.Println("\rChild nodes:")
	for _, child := range node.Children {
		fmt.Println("\r", child.Name)
	}
	fmt.Println("\rMarkdown files:")
	for _, markdown := range node.Markdowns {
		fmt.Println("\r", markdown.Filename)
	}
}

func (mmf *MetadataNoteWolfyFileHandle) FindNode(name string) *Node {
	activeWorkspaceName := mmf.ActiveWorkspace
	var activeWorkspace *Node
	for _, workspace := range mmf.Workspaces {
		if workspace.Name == activeWorkspaceName {
			activeWorkspace = workspace
		}
	}
	if activeWorkspace == nil {
		return nil
	}

	queue := NewQueue[*Node]()
	queue.Add(activeWorkspace)

	var node *Node
	for queue.Len() > 0 {
		currentNode := queue.Drop()
		if currentNode.Name == name {
			node = currentNode
			break
		}
		for _, child := range currentNode.Children {
			queue.Add(child)
		}
	}

	return node
}

func (mmf *MetadataNoteWolfyFileHandle) FindParentNode(name string) *Node {
	activeWorkspaceName := mmf.ActiveWorkspace
	var activeWorkspace *Node
	for _, workspace := range mmf.Workspaces {
		if workspace.Name == activeWorkspaceName {
			activeWorkspace = workspace
		}
	}
	if activeWorkspace == nil {
		return nil
	}

	queue := NewQueue[*Node]()
	queue.Add(activeWorkspace)

	var parentNode *Node
	for queue.Len() > 0 {
		currentNode := queue.Drop()
		for _, childNode := range currentNode.Children {
			if childNode.Name == name {
				parentNode = currentNode
				break
			}
			queue.Add(childNode)
		}
	}

	return parentNode
}

func (mmf *MetadataNoteWolfyFileHandle) load() error {
	file, err := mmf.getMetadataFile()
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	if fileInfo.Size() == 0 {
		return nil
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(mmf); err != nil {
		return err
	}
	return nil
}

func (mmf *MetadataNoteWolfyFileHandle) getMetadataFile() (*os.File, error) {
	metadataFilePath := mmf.Config.MetadataFilePath

	file, err := os.OpenFile(metadataFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return file, nil
}

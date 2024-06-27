//go:build unit_test

package structure_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/RaphSku/notewolfy/internal/structure"
	"github.com/RaphSku/notewolfy/internal/utility"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func CleanUpFile(filePath string) {
	err := os.Remove(filePath)
	if err != nil {
		fmt.Println("Could not remove metadata file! Please clean it up yourself!")
		os.Exit(1)
	}
}

func captureStdOutput(f func()) (string, error) {
	originalStdOut := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}
	os.Stdout = w

	outputC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outputC <- buf.String()
	}()

	f()
	w.Close()

	os.Stdout = originalStdOut
	out := <-outputC

	return out, nil
}

func TestNewMetadataNoteWolfyFileHandle(t *testing.T) {
	t.Parallel()

	uuid := uuid.New().String()
	metadataFilePath := fmt.Sprintf("./.notewolfy-%s.json", uuid)
	config := &structure.Config{
		MetadataFilePath: metadataFilePath,
	}
	defer CleanUpFile(metadataFilePath)

	var expMmf structure.MetadataNoteWolfyFileHandle
	expMmf.Config = config

	mmf1, err := structure.NewMetadataNoteWolfyFileHandle(config)
	assert.NoError(t, err)
	mmf2, err := structure.NewMetadataNoteWolfyFileHandle(config)
	assert.NoError(t, err)
	assert.Equal(t, mmf1, mmf2)

	assert.Equal(t, expMmf.Config, mmf1.Config)
	assert.Equal(t, expMmf.Workspaces, mmf1.Workspaces)
	assert.Equal(t, expMmf.ActiveWorkspace, mmf1.ActiveWorkspace)
}

func TestCreateNewWorkspace(t *testing.T) {
	t.Parallel()

	uuid := uuid.New().String()
	metadataFilePath := fmt.Sprintf("./.notewolfy-%s.json", uuid)
	config := &structure.Config{
		MetadataFilePath: metadataFilePath,
	}
	defer CleanUpFile(metadataFilePath)

	mmf, err := structure.NewMetadataNoteWolfyFileHandle(config)
	assert.NoError(t, err)

	expRelativeWorkspacePath := "~/tmp"
	err = mmf.AddNewWorkspace("test", expRelativeWorkspacePath)
	assert.NoError(t, err)

	expWorkspacePath, err := utility.ExpandRelativePaths(expRelativeWorkspacePath)
	assert.NoError(t, err)
	expWorkspaceNode := &structure.Node{
		Name: "test",
		Path: expWorkspacePath,
	}
	expWorkspaceList := []*structure.Node{expWorkspaceNode}
	assert.Equal(t, expWorkspaceList, mmf.Workspaces)
	assert.Equal(t, expWorkspaceList[0].Path, mmf.Workspaces[0].Path)
}

func TestDoesWorkspaceExist(t *testing.T) {
	t.Parallel()

	uuid := uuid.New().String()
	metadataFilePath := fmt.Sprintf("./.notewolfy-%s.json", uuid)
	config := &structure.Config{
		MetadataFilePath: metadataFilePath,
	}
	defer CleanUpFile(metadataFilePath)

	mmf, err := structure.NewMetadataNoteWolfyFileHandle(config)
	assert.NoError(t, err)

	expExists := false
	actExists := mmf.DoesWorkspaceExist("test")
	assert.Equal(t, expExists, actExists)

	err = mmf.AddNewWorkspace("test", "~/tmp")
	assert.NoError(t, err)
	expExists = true
	actExists = mmf.DoesWorkspaceExist("test")
	assert.Equal(t, expExists, actExists)
}

func TestSave(t *testing.T) {
	t.Parallel()

	uuid := uuid.New().String()
	metadataFilePath := fmt.Sprintf("./.notewolfy-%s.json", uuid)
	config := &structure.Config{
		MetadataFilePath: metadataFilePath,
	}
	defer CleanUpFile(metadataFilePath)

	mmf, err := structure.NewMetadataNoteWolfyFileHandle(config)
	assert.NoError(t, err)

	err = mmf.AddNewWorkspace("testA", "~/tmp")
	assert.NoError(t, err)
	err = mmf.AddNewWorkspace("testB", "~/A/B")
	assert.NoError(t, err)

	err = mmf.Save()
	assert.NoError(t, err)

	file, err := os.Open(metadataFilePath)
	assert.NoError(t, err)
	decoder := json.NewDecoder(file)
	var actMmf structure.MetadataNoteWolfyFileHandle
	err = decoder.Decode(&actMmf)
	assert.NoError(t, err)

	for i := range mmf.Workspaces {
		assert.Equal(t, mmf.Workspaces[i].Name, actMmf.Workspaces[i].Name)
	}
	assert.Equal(t, mmf.ActiveWorkspace, actMmf.ActiveWorkspace)
	assert.Equal(t, mmf.ActiveNode, actMmf.ActiveNode)
}

func TestAddChildToNodeFailure(t *testing.T) {
	t.Parallel()

	uuid := uuid.New().String()
	metadataFilePath := fmt.Sprintf("./.notewolfy-%s.json", uuid)
	config := &structure.Config{
		MetadataFilePath: metadataFilePath,
	}
	defer CleanUpFile(metadataFilePath)

	mmf, err := structure.NewMetadataNoteWolfyFileHandle(config)
	assert.NoError(t, err)

	childNode := &structure.Node{
		Name: "B",
		Path: "/B",
	}

	err = mmf.AddChild(childNode)
	if assert.Error(t, err) {
		expError := errors.New("no active node, seems like you have not created a workspace yet")
		assert.Equal(t, expError, err)
	}

	node := &structure.Node{
		Name: "A",
		Path: "/A",
	}
	mmf.ActiveNode = node.Name
	mmf.ActiveWorkspace = node.Name
	mmf.Workspaces = append(mmf.Workspaces, node)
	err = mmf.AddChild(childNode)
	if assert.Error(t, err) {
		expError := errors.New("child's parent path and parentPath do not match")
		assert.Equal(t, expError, err)
	}
}

func TestAddChildToNodeSuccess(t *testing.T) {
	t.Parallel()

	uuid := uuid.New().String()
	metadataFilePath := fmt.Sprintf("./.notewolfy-%s.json", uuid)
	config := &structure.Config{
		MetadataFilePath: metadataFilePath,
	}
	defer CleanUpFile(metadataFilePath)

	mmf, err := structure.NewMetadataNoteWolfyFileHandle(config)
	assert.NoError(t, err)

	var markdowns []*structure.Markdown
	var children []*structure.Node
	node := &structure.Node{
		Name:      "A",
		Path:      "~/A",
		Markdowns: markdowns,
		Children:  children,
	}
	mmf.Workspaces = append(mmf.Workspaces, node)
	mmf.ActiveNode = node.Name
	mmf.ActiveWorkspace = node.Name

	childNode := &structure.Node{
		Name:      "B",
		Path:      "~/A/B",
		Markdowns: markdowns,
		Children:  children,
	}
	err = mmf.AddChild(childNode)
	assert.NoError(t, err)
	actNode := mmf.FindNode(node.Name)
	assert.Equal(t, actNode.Children[0], childNode)
	mmf.Save()
}

func TestDeleteChildByIndex(t *testing.T) {
	t.Parallel()

	uuid := uuid.New().String()
	metadataFilePath := fmt.Sprintf("./.notewolfy-%s.json", uuid)
	config := &structure.Config{
		MetadataFilePath: metadataFilePath,
	}
	defer CleanUpFile(metadataFilePath)

	mmf, err := structure.NewMetadataNoteWolfyFileHandle(config)
	assert.NoError(t, err)

	nodeA := &structure.Node{
		Name: "A",
		Path: "~/test/A",
	}
	nodeB := &structure.Node{
		Name: "B",
		Path: "~/test/B",
	}
	nodeC := &structure.Node{
		Name: "C",
		Path: "~/test/C",
	}
	workspaceNode := &structure.Node{
		Name:     "Test",
		Path:     "~/test",
		Children: []*structure.Node{nodeA, nodeB, nodeC},
	}

	err = mmf.DeleteChildByIndex(0)
	if assert.Error(t, err) {
		expErr := errors.New("no active node, seems like you have not created a workspace yet")
		assert.Equal(t, expErr, err)
	}

	mmf.ActiveNode = workspaceNode.Name
	mmf.ActiveWorkspace = workspaceNode.Name
	mmf.Workspaces = append(mmf.Workspaces, workspaceNode)

	err = mmf.DeleteChildByIndex(-1)
	if assert.Error(t, err) {
		expErr := errors.New("index is out of range, check that the child at this index exists")
		assert.Equal(t, expErr, err)
	}

	err = mmf.DeleteChildByIndex(4)
	if assert.Error(t, err) {
		expErr := errors.New("index is out of range, check that the child at this index exists")
		assert.Equal(t, expErr, err)
	}

	err = mmf.DeleteChildByIndex(1)
	expLength := 2
	assert.NoError(t, err)
	assert.Equal(t, expLength, len(mmf.Workspaces[0].Children))
	assert.True(t, reflect.DeepEqual(nodeA, mmf.Workspaces[0].Children[0]))
	assert.True(t, reflect.DeepEqual(nodeC, mmf.Workspaces[0].Children[1]))
}

func TestAddMarkdown(t *testing.T) {
	t.Parallel()

	uuid := uuid.New().String()
	metadataFilePath := fmt.Sprintf("./.notewolfy-%s.json", uuid)
	config := &structure.Config{
		MetadataFilePath: metadataFilePath,
	}
	defer CleanUpFile(metadataFilePath)

	mmf, err := structure.NewMetadataNoteWolfyFileHandle(config)
	assert.NoError(t, err)

	var markdowns []*structure.Markdown
	var children []*structure.Node
	node := &structure.Node{
		Name:      "test",
		Path:      "/A",
		Markdowns: markdowns,
		Children:  children,
	}
	mmf.Workspaces = append(mmf.Workspaces, node)
	mmf.ActiveNode = node.Name
	mmf.ActiveWorkspace = node.Name

	markdown := &structure.Markdown{
		Filename: "test.md",
	}
	err = mmf.AddMarkdown(markdown)
	assert.NoError(t, err)
	actNode := mmf.FindNode(node.Name)
	assert.Equal(t, actNode.Markdowns[0].Filename, markdown.Filename)
}

func TestDeleteMarkdown(t *testing.T) {
	t.Parallel()

	uuid := uuid.New().String()
	metadataFilePath := fmt.Sprintf("./.notewolfy-%s.json", uuid)
	config := &structure.Config{
		MetadataFilePath: metadataFilePath,
	}
	defer CleanUpFile(metadataFilePath)

	mmf, err := structure.NewMetadataNoteWolfyFileHandle(config)
	assert.NoError(t, err)

	var markdowns []*structure.Markdown
	var children []*structure.Node
	node := &structure.Node{
		Name:      "test",
		Path:      "/A",
		Markdowns: markdowns,
		Children:  children,
	}
	mmf.ActiveNode = node.Name
	mmf.ActiveWorkspace = node.Name
	mmf.Workspaces = append(mmf.Workspaces, node)

	markdownName := "test"
	markdown := &structure.Markdown{
		Filename: fmt.Sprintf("%s.md", markdownName),
	}
	err = mmf.AddMarkdown(markdown)
	assert.NoError(t, err)

	err = mmf.DeleteMarkdown(markdownName)
	assert.NoError(t, err)
	actNode := mmf.FindNode(node.Name)
	assert.Empty(t, actNode.Markdowns)
}

func TestListWorkspaces(t *testing.T) {
	uuid := uuid.New().String()
	metadataFilePath := fmt.Sprintf("./.notewolfy-%s.json", uuid)
	config := &structure.Config{
		MetadataFilePath: metadataFilePath,
	}
	defer CleanUpFile(metadataFilePath)

	mmf, err := structure.NewMetadataNoteWolfyFileHandle(config)
	assert.NoError(t, err)

	var workspaceNodes []*structure.Node
	var markdowns []*structure.Markdown
	var children []*structure.Node
	workspaceNodes = append(workspaceNodes, &structure.Node{
		Name:      "A",
		Path:      "./tmp/A",
		Markdowns: markdowns,
		Children:  children,
	})
	workspaceNodes = append(workspaceNodes, &structure.Node{
		Name:      "B",
		Path:      "./tmp/B",
		Markdowns: markdowns,
		Children:  children,
	})
	mmf.Workspaces = workspaceNodes

	actOutput, err := captureStdOutput(func() {
		mmf.ListWorkspaces()
	})
	assert.NoError(t, err)

	expOutput := "\r\nWorkspace NameWorkspace Path\r\n-----------------------------\n\rA             ./tmp/A       \n\rB             ./tmp/B       \n"
	assert.Equal(t, expOutput, actOutput)
}

func TestListResourcesOnNode(t *testing.T) {
	uuid := uuid.New().String()
	metadataFilePath := fmt.Sprintf("./.notewolfy-%s.json", uuid)
	config := &structure.Config{
		MetadataFilePath: metadataFilePath,
	}
	defer CleanUpFile(metadataFilePath)

	mmf, err := structure.NewMetadataNoteWolfyFileHandle(config)
	assert.NoError(t, err)

	var markdowns []*structure.Markdown
	markdown := &structure.Markdown{
		Filename: "test",
	}
	markdowns = append(markdowns, markdown)

	var children []*structure.Node
	childNode := &structure.Node{
		Name: "childTest",
		Path: "/A/B",
	}
	children = append(children, childNode)

	node := &structure.Node{
		Name:      "test",
		Path:      "/A",
		Markdowns: markdowns,
		Children:  children,
	}

	mmf.ActiveNode = node.Name
	mmf.ActiveWorkspace = node.Name
	mmf.Workspaces = append(mmf.Workspaces, node)

	actOutput, err := captureStdOutput(func() {
		mmf.ListResourcesOnNode(node)
	})
	assert.NoError(t, err)

	expOutput := "\r\nYou are on node:  test\n\rChild nodes:\n\r childTest\n\rMarkdown files:\n\r test\n"
	assert.Equal(t, expOutput, actOutput)
}

func TestFindNode(t *testing.T) {
	t.Parallel()

	uuid := uuid.New().String()
	metadataFilePath := fmt.Sprintf("./.notewolfy-%s.json", uuid)
	config := &structure.Config{
		MetadataFilePath: metadataFilePath,
	}
	defer CleanUpFile(metadataFilePath)

	mmf, err := structure.NewMetadataNoteWolfyFileHandle(config)
	assert.NoError(t, err)

	var childrenOfB []*structure.Node
	nodeD := &structure.Node{
		Name: "D",
		Path: "/A/B/D",
	}
	childrenOfB = append(childrenOfB, nodeD)
	nodeB := &structure.Node{
		Name:     "B",
		Path:     "/A/B",
		Children: childrenOfB,
	}
	nodeC := &structure.Node{
		Name: "C",
		Path: "/A/C",
	}
	var children []*structure.Node
	children = append(children, nodeB)
	children = append(children, nodeC)
	nodeA := &structure.Node{
		Name:     "A",
		Path:     "/A",
		Children: children,
	}
	mmf.Workspaces = append(mmf.Workspaces, nodeA)
	mmf.ActiveWorkspace = nodeA.Name
	mmf.ActiveNode = nodeA.Name

	actNode := mmf.FindNode("D")
	assert.True(t, reflect.DeepEqual(nodeD, actNode))
}

func TestFindParentNode(t *testing.T) {
	t.Parallel()

	uuid := uuid.New().String()
	metadataFilePath := fmt.Sprintf("./.notewolfy-%s.json", uuid)
	config := &structure.Config{
		MetadataFilePath: metadataFilePath,
	}
	defer CleanUpFile(metadataFilePath)

	mmf, err := structure.NewMetadataNoteWolfyFileHandle(config)
	assert.NoError(t, err)

	var childrenOfB []*structure.Node
	nodeD := &structure.Node{
		Name: "D",
		Path: "/A/B/D",
	}
	childrenOfB = append(childrenOfB, nodeD)
	nodeB := &structure.Node{
		Name:     "B",
		Path:     "/A/B",
		Children: childrenOfB,
	}
	nodeC := &structure.Node{
		Name: "C",
		Path: "/A/C",
	}
	var children []*structure.Node
	children = append(children, nodeB)
	children = append(children, nodeC)
	nodeA := &structure.Node{
		Name:     "A",
		Path:     "/A",
		Children: children,
	}
	mmf.Workspaces = append(mmf.Workspaces, nodeA)
	mmf.ActiveWorkspace = nodeA.Name
	mmf.ActiveNode = nodeA.Name

	actParentNode := mmf.FindParentNode("D")
	assert.True(t, reflect.DeepEqual(nodeB, actParentNode))
}

func TestDecodingWhileLoading(t *testing.T) {
	t.Parallel()

	uuid := uuid.New().String()
	metadataFilePath := fmt.Sprintf("./.notewolfy-%s.json", uuid)
	config := &structure.Config{
		MetadataFilePath: metadataFilePath,
	}
	defer CleanUpFile(metadataFilePath)

	mmf, err := structure.NewMetadataNoteWolfyFileHandle(config)
	assert.NoError(t, err)

	activeNode := &structure.Node{
		Name: "Active",
		Path: "/A",
	}
	mmf.ActiveWorkspace = activeNode.Name
	mmf.ActiveNode = activeNode.Name
	mmf.Workspaces = append(mmf.Workspaces, activeNode)

	err = mmf.Save()
	assert.NoError(t, err)

	newConfig := &structure.Config{
		MetadataFilePath: metadataFilePath,
	}
	newMmf, err := structure.NewMetadataNoteWolfyFileHandle(newConfig)
	assert.NoError(t, err)

	assert.True(t, reflect.DeepEqual(activeNode, newMmf.Workspaces[0]))
}

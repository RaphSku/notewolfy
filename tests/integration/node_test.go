//go:build integration_test

package integration_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/RaphSku/notewolfy/internal/commands"
	"github.com/RaphSku/notewolfy/internal/structure"
	"github.com/RaphSku/notewolfy/internal/utility"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func createUniquePath(path string) string {
	uuid := uuid.New().String()

	return fmt.Sprintf("%s-%s", path, uuid)
}

func CleanUpFile(filePath string) {
	err := os.Remove(filePath)
	if err != nil {
		fmt.Println("Could not remove metadata file! Please clean it up yourself!")
		os.Exit(1)
	}
}

func TestNodeCreatingAndDeleting(t *testing.T) {
	// Scenario:
	// 1. Prepare workspace
	// 2. Create test node
	// 3. Create another node
	// 4. Use command `goto` to move to test node
	// 5. Use command `goback` to move back to the workspace node
	// 6. Delete test node
	t.Parallel()

	metadataFilePath := createUniquePath("./.notewolfy")
	config := &structure.Config{
		MetadataFilePath: metadataFilePath,
	}
	defer CleanUpFile(metadataFilePath)

	mmf, err := structure.NewMetadataNoteWolfyFileHandle(config)
	assert.NoError(t, err)

	workspacePath := createUniquePath("./tmp")

	// 1. Prepare workspace
	workspacePath, err = utility.ExpandRelativePaths(workspacePath)
	assert.NoError(t, err)
	err = os.Mkdir(workspacePath, os.ModePerm)
	defer os.RemoveAll(workspacePath)
	assert.NoError(t, err)
	workspaceNode := &structure.Node{
		Name: "Workspace",
		Path: workspacePath,
	}
	mmf.Workspaces = append(mmf.Workspaces, workspaceNode)
	mmf.ActiveNode = workspaceNode.Name
	mmf.ActiveWorkspace = workspaceNode.Name
	err = mmf.Save()
	assert.NoError(t, err)

	// 2. Create test node
	testNodeName := "Test"
	statement := fmt.Sprintf("create node %s", testNodeName)
	commands.MatchStatementToCommand(mmf, statement)
	testNodePath := filepath.Join(workspacePath, testNodeName)
	assert.DirExists(t, testNodePath)
	assert.Equal(t, testNodeName, mmf.Workspaces[0].Children[0].Name)

	// 3. Create another node
	nodeName := "A"
	statement = fmt.Sprintf("create node %s", nodeName)
	commands.MatchStatementToCommand(mmf, statement)
	nodePath := filepath.Join(workspacePath, testNodeName)
	defer os.Remove(nodePath)
	assert.DirExists(t, nodePath)

	// 4. Use command `goto` to move to test node
	statement = fmt.Sprintf("goto %s", testNodeName)
	commands.MatchStatementToCommand(mmf, statement)
	assert.Equal(t, testNodeName, mmf.ActiveNode)

	// 5. Use command `goback` to move back to the workspace node
	statement = "goback"
	commands.MatchStatementToCommand(mmf, statement)
	assert.Equal(t, workspaceNode.Name, mmf.ActiveNode)

	// 6. Delete test node
	statement = fmt.Sprintf("delete node %s", testNodeName)
	commands.MatchStatementToCommand(mmf, statement)
	if _, err = os.Stat(testNodePath); err == nil {
		os.Remove(testNodePath)
	}
	assert.NoDirExists(t, testNodePath)
	assert.Equal(t, 1, len(mmf.Workspaces[0].Children))
	assert.Equal(t, nodeName, mmf.Workspaces[0].Children[0].Name)
}

//go:build unit_test

package commands_test

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/RaphSku/notewolfy/internal/commands"
	"github.com/RaphSku/notewolfy/internal/structure"
	"github.com/RaphSku/notewolfy/internal/utility"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func fileOrDirectoryExists(path string) (bool, error) {
	expandedPath, err := utility.ExpandRelativePaths(path)
	if err != nil {
		return false, err
	}
	_, err = os.Stat(expandedPath)
	if err != nil {
		return false, err
	}
	return true, nil
}

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

func TestMatchStatementToCreateWorkspaceCommand(t *testing.T) {
	t.Parallel()

	firstTestCasePath := createUniquePath("./tmp")
	secondTestCasePath := createUniquePath("~/tmp")
	thirdTestCasePath := createUniquePath("./tmp")
	fourthTestCasePath := createUniquePath("./tmp")

	tests := map[string]struct {
		statement string
		path      string
		want      bool
	}{
		"simple create workspace command with relative path": {statement: fmt.Sprintf("create workspace %s %s", "test", firstTestCasePath), path: firstTestCasePath, want: true},
		"simple create workspace command with absolute path": {statement: fmt.Sprintf("create workspace %s %s", "test", secondTestCasePath), path: secondTestCasePath, want: true},
		"error create workspace command":                     {statement: fmt.Sprintf("dgkhs create workspace %s %s", "test", thirdTestCasePath), path: thirdTestCasePath, want: false},
		"error create workspace command that almost matches": {statement: fmt.Sprintf("create workspaces %s %s", "test", fourthTestCasePath), path: fourthTestCasePath, want: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			metadataFilePath := createUniquePath("./.notewolfy")
			config := &structure.Config{
				MetadataFilePath: metadataFilePath,
			}
			defer CleanUpFile(metadataFilePath)

			mmf, err := structure.NewMetadataNoteWolfyFileHandle(config)
			assert.NoError(t, err)

			statement := tc.statement
			commands.MatchStatementToCommand(mmf, statement)
			exists, err := fileOrDirectoryExists(tc.path)
			if tc.want {
				expandedPath, err := utility.ExpandRelativePaths(tc.path)
				assert.NoError(t, err)
				defer os.Remove(expandedPath)

				assert.NoError(t, err)
				assert.True(t, exists)

				return
			}
			if assert.Error(t, err) {
				if _, ok := err.(*fs.PathError); ok {
					assert.True(t, ok)
				}
			}
			assert.False(t, exists)
		})
	}
}

func TestMatchStatementToDeleteWorkspaceCommand(t *testing.T) {
	t.Parallel()

	firstTestCasePath := createUniquePath("./tmp")
	secondTestCasePath := createUniquePath("~/tmp")
	thirdTestCasePath := createUniquePath("./tmp")
	fourthTestCasePath := createUniquePath("./tmp")

	tests := map[string]struct {
		statement string
		path      string
		want      bool
	}{
		"simple delete workspace command with relative path": {
			statement: fmt.Sprintf("delete workspace %s %s", "test", firstTestCasePath),
			path:      firstTestCasePath,
			want:      true,
		},
		"simple delete workspace command with absolute path": {
			statement: fmt.Sprintf("delete workspace %s %s", "test", secondTestCasePath),
			path:      secondTestCasePath,
			want:      true,
		},
		"error delete workspace command": {
			statement: fmt.Sprintf("dgkhs delete workspace %s %s", "test", thirdTestCasePath),
			path:      thirdTestCasePath,
			want:      false,
		},
		"error delete workspace command that almost matches": {
			statement: fmt.Sprintf("delete workspaces %s %s", "test", fourthTestCasePath),
			path:      fourthTestCasePath,
			want:      false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			metadataFilePath := createUniquePath("./.notewolfy")
			config := &structure.Config{
				MetadataFilePath: metadataFilePath,
			}
			defer CleanUpFile(metadataFilePath)

			mmf, err := structure.NewMetadataNoteWolfyFileHandle(config)
			assert.NoError(t, err)

			// Need to prepare the workspace
			workspacePath, err := utility.ExpandRelativePaths(tc.path)
			assert.NoError(t, err)
			err = os.Mkdir(workspacePath, os.ModePerm)
			assert.NoError(t, err)
			workspaceNode := &structure.Node{
				Name: "test",
				Path: workspacePath,
			}
			mmf.Workspaces = append(mmf.Workspaces, workspaceNode)
			mmf.ActiveWorkspace = workspaceNode.Name
			mmf.ActiveNode = workspaceNode.Name
			err = mmf.Save()
			assert.NoError(t, err)

			defer os.Remove(workspaceNode.Path)

			commands.MatchStatementToCommand(mmf, tc.statement)
			_, err = fileOrDirectoryExists(tc.path)
			if tc.want {
				if assert.Error(t, err) {
					if _, ok := err.(*fs.PathError); ok {
						assert.True(t, ok)
					}
				}
				assert.Empty(t, mmf.ActiveNode)
				assert.Empty(t, mmf.ActiveWorkspace)
				assert.Equal(t, 0, len(mmf.Workspaces))

				return
			}
			assert.DirExists(t, tc.path)
		})
	}
}

func TestMatchStatementToCreateNodeCommand(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		statement     string
		workspacePath string
		nodeName      string
		want          bool
	}{
		"simple create node command": {
			statement:     "create node A",
			workspacePath: createUniquePath("./tmp"),
			nodeName:      "A",
			want:          true,
		},
		"error create node command": {
			statement:     "dgkhs create node test",
			workspacePath: createUniquePath("./tmp"),
			nodeName:      "test",
			want:          false,
		},
		"error create node command that almost matches": {
			statement:     "create nodes test",
			workspacePath: createUniquePath("./tmp"),
			nodeName:      "test",
			want:          false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			metadataFilePath := createUniquePath("./.notewolfy")
			config := &structure.Config{
				MetadataFilePath: metadataFilePath,
			}
			defer CleanUpFile(metadataFilePath)

			mmf, err := structure.NewMetadataNoteWolfyFileHandle(config)
			assert.NoError(t, err)

			// We need to prepare the workspace before we can create a node
			workspacePath, err := utility.ExpandRelativePaths(tc.workspacePath)
			assert.NoError(t, err)
			err = os.Mkdir(workspacePath, os.ModePerm)
			assert.NoError(t, err)
			workspaceNode := &structure.Node{
				Name: "Workspace",
				Path: workspacePath,
			}
			mmf.Workspaces = append(mmf.Workspaces, workspaceNode)
			mmf.ActiveWorkspace = workspaceNode.Name
			mmf.ActiveNode = workspaceNode.Name
			err = mmf.Save()
			assert.NoError(t, err)

			statement := tc.statement
			commands.MatchStatementToCommand(mmf, statement)
			nodePath := filepath.Join(workspacePath, tc.nodeName)
			exists, err := fileOrDirectoryExists(nodePath)
			if tc.want {
				os.Remove(nodePath)
				os.Remove(workspacePath)

				assert.NoError(t, err)
				assert.True(t, exists)

				return
			}
			os.Remove(workspacePath)
			if assert.Error(t, err) {
				if _, ok := err.(*fs.PathError); ok {
					assert.True(t, ok)
				}
			}
			assert.False(t, exists)
		})
	}
}

func TestMatchStatementToDeleteNodeCommand(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		statement     string
		workspacePath string
		nodeName      string
		want          bool
	}{
		"simple delete node command": {
			statement:     "delete node A",
			workspacePath: createUniquePath("./tmp"),
			nodeName:      "A",
			want:          true,
		},
		"error delete node command": {
			statement:     "dgkhs delete node A",
			workspacePath: createUniquePath("./tmp"),
			nodeName:      "A",
			want:          false,
		},
		"error delete node command that almost matches": {
			statement:     "delete nodes A",
			workspacePath: createUniquePath("./tmp"),
			nodeName:      "A",
			want:          false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			metadataFilePath := createUniquePath("./.notewolfy")
			config := &structure.Config{
				MetadataFilePath: metadataFilePath,
			}
			defer CleanUpFile(metadataFilePath)

			mmf, err := structure.NewMetadataNoteWolfyFileHandle(config)
			assert.NoError(t, err)

			// We need to prepare the workspace before we can create & delete a node
			workspacePath, err := utility.ExpandRelativePaths(tc.workspacePath)
			assert.NoError(t, err)
			err = os.Mkdir(workspacePath, os.ModePerm)
			assert.NoError(t, err)
			workspaceNode := &structure.Node{
				Name: "Workspace",
				Path: workspacePath,
			}
			mmf.Workspaces = append(mmf.Workspaces, workspaceNode)
			mmf.ActiveWorkspace = workspaceNode.Name
			mmf.ActiveNode = workspaceNode.Name
			err = mmf.Save()
			assert.NoError(t, err)

			defer os.Remove(workspacePath)

			// Let's create a node that can be deleted
			nodePath := filepath.Join(workspacePath, tc.nodeName)
			err = os.Mkdir(nodePath, os.ModePerm)
			assert.NoError(t, err)
			node := &structure.Node{
				Name: tc.nodeName,
				Path: nodePath,
			}
			mmf.Workspaces[0].Children = append(mmf.Workspaces[0].Children, node)
			err = mmf.Save()
			assert.NoError(t, err)

			defer os.Remove(nodePath)

			statement := tc.statement
			commands.MatchStatementToCommand(mmf, statement)
			_, err = fileOrDirectoryExists(nodePath)
			if tc.want {
				if assert.Error(t, err) {
					if _, ok := err.(*fs.PathError); ok {
						assert.True(t, ok)
					}
				}
				assert.Equal(t, 0, len(mmf.Workspaces[0].Children))
				return
			}
			assert.DirExists(t, nodePath)
		})
	}
}

func TestMatchStatementToCreateMarkdownFile(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		statement     string
		workspacePath string
		markdownName  string
		want          bool
	}{
		"simple create markdown command": {
			statement:     "create md example",
			workspacePath: createUniquePath("./tmp"),
			markdownName:  "example",
			want:          true,
		},
		"error create markdown command": {
			statement:     "dgkhs create markdown example some",
			workspacePath: createUniquePath("./tmp"),
			markdownName:  "example",
			want:          false,
		},
		"error create markdown command that almost matches": {
			statement:     "create markdowns example",
			workspacePath: createUniquePath("./tmp"),
			markdownName:  "example",
			want:          false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			metadataFilePath := createUniquePath("./.notewolfy")
			config := &structure.Config{
				MetadataFilePath: metadataFilePath,
			}
			defer CleanUpFile(metadataFilePath)

			mmf, err := structure.NewMetadataNoteWolfyFileHandle(config)
			assert.NoError(t, err)

			// We need to prepare the workspace
			workspacePath, err := utility.ExpandRelativePaths(tc.workspacePath)
			assert.NoError(t, err)
			err = os.Mkdir(workspacePath, os.ModePerm)
			assert.NoError(t, err)
			workspaceNode := &structure.Node{
				Name: "Workspace",
				Path: workspacePath,
			}
			mmf.Workspaces = append(mmf.Workspaces, workspaceNode)
			mmf.ActiveWorkspace = workspaceNode.Name
			mmf.ActiveNode = workspaceNode.Name
			err = mmf.Save()
			assert.NoError(t, err)

			statement := tc.statement
			commands.MatchStatementToCommand(mmf, statement)
			markdownPath := filepath.Join(workspacePath, tc.markdownName+".md")
			exists, err := fileOrDirectoryExists(markdownPath)
			if tc.want {
				os.Remove(markdownPath)
				os.Remove(workspacePath)

				assert.NoError(t, err)
				assert.True(t, exists)

				return
			}
			os.Remove(workspacePath)
			if assert.Error(t, err) {
				if _, ok := err.(*fs.PathError); ok {
					assert.True(t, ok)
				}
			}
			assert.False(t, exists)
		})
	}
}

func TestMatchStatementToDeleteMarkdownFile(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		statement     string
		workspacePath string
		markdownName  string
		want          bool
	}{
		"simple delete markdown command": {
			statement:     "delete md example",
			workspacePath: createUniquePath("./tmp"),
			markdownName:  "example",
			want:          true,
		},
		"error delete markdown command": {
			statement:     "dgkhs delete md example some",
			workspacePath: createUniquePath("./tmp"),
			markdownName:  "example",
			want:          false,
		},
		"error delete markdown command that almost matches": {
			statement:     "delete markdown example",
			workspacePath: createUniquePath("./tmp"),
			markdownName:  "example",
			want:          false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			metadataFilePath := createUniquePath("./.notewolfy")
			config := &structure.Config{
				MetadataFilePath: metadataFilePath,
			}
			defer CleanUpFile(metadataFilePath)

			mmf, err := structure.NewMetadataNoteWolfyFileHandle(config)
			assert.NoError(t, err)

			// We need to prepare the workspace
			workspacePath, err := utility.ExpandRelativePaths(tc.workspacePath)
			assert.NoError(t, err)
			err = os.Mkdir(workspacePath, os.ModePerm)
			assert.NoError(t, err)
			workspaceNode := &structure.Node{
				Name: "Workspace",
				Path: workspacePath,
			}
			mmf.Workspaces = append(mmf.Workspaces, workspaceNode)
			mmf.ActiveWorkspace = workspaceNode.Name
			mmf.ActiveNode = workspaceNode.Name
			err = mmf.Save()
			assert.NoError(t, err)

			defer os.Remove(workspacePath)

			// We need to create a Markdown file that we can delete
			markdown := &structure.Markdown{
				Filename: tc.markdownName + ".md",
			}
			markdownPath := filepath.Join(workspacePath, tc.markdownName+".md")
			file, err := os.Create(markdownPath)
			assert.NoError(t, err)
			defer os.Remove(markdownPath)
			defer file.Close()
			mmf.Workspaces[0].Markdowns = append(mmf.Workspaces[0].Markdowns, markdown)
			err = mmf.Save()
			assert.NoError(t, err)

			statement := tc.statement
			commands.MatchStatementToCommand(mmf, statement)
			exists, err := fileOrDirectoryExists(markdownPath)
			if tc.want {
				if assert.Error(t, err) {
					if _, ok := err.(*fs.PathError); ok {
						assert.True(t, ok)
					}
				}
				assert.False(t, exists)

				return
			}
			assert.NoError(t, err)
			assert.True(t, exists)
		})
	}
}

func TestMatchStatementToGotoNode(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		statement     string
		workspacePath string
		nodeName      string
		want          bool
	}{
		"simple goto command": {
			statement:     "goto example",
			workspacePath: createUniquePath("./tmp"),
			nodeName:      "example",
			want:          true,
		},
		"error goto command": {
			statement:     "gotos example",
			workspacePath: createUniquePath("./tmp"),
			nodeName:      "example",
			want:          false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			metadataFilePath := createUniquePath("./.notewolfy")
			config := &structure.Config{
				MetadataFilePath: metadataFilePath,
			}
			defer CleanUpFile(metadataFilePath)

			mmf, err := structure.NewMetadataNoteWolfyFileHandle(config)
			assert.NoError(t, err)

			// We need to prepare the workspace
			workspacePath, err := utility.ExpandRelativePaths(tc.workspacePath)
			assert.NoError(t, err)
			err = os.Mkdir(workspacePath, os.ModePerm)
			assert.NoError(t, err)
			workspaceNode := &structure.Node{
				Name: "Workspace",
				Path: workspacePath,
			}
			mmf.Workspaces = append(mmf.Workspaces, workspaceNode)
			mmf.ActiveWorkspace = workspaceNode.Name
			mmf.ActiveNode = workspaceNode.Name
			err = mmf.Save()
			assert.NoError(t, err)

			// We need to create a node
			nodePath := filepath.Join(tc.workspacePath, tc.nodeName)
			err = os.Mkdir(nodePath, os.ModePerm)
			assert.NoError(t, err)
			node := &structure.Node{
				Name: tc.nodeName,
				Path: nodePath,
			}
			mmf.Workspaces[0].Children = append(mmf.Workspaces[0].Children, node)
			err = mmf.Save()
			assert.NoError(t, err)

			defer os.Remove(workspacePath)
			defer os.Remove(nodePath)

			statement := tc.statement
			commands.MatchStatementToCommand(mmf, statement)
			if tc.want {
				assert.Equal(t, node.Name, mmf.ActiveNode)
				actNode := mmf.FindNode(node.Name)
				assert.True(t, reflect.DeepEqual(node, actNode))

				return
			}
			assert.Equal(t, workspaceNode.Name, mmf.ActiveNode)
		})
	}
}

func TestMatchStatementToGoBack(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		statement     string
		workspacePath string
		nodeName      string
		want          bool
	}{
		"simple goback command": {
			statement:     "goback",
			workspacePath: createUniquePath("./tmp"),
			nodeName:      "example",
			want:          true,
		},
		"error goback command": {
			statement:     "gobacks",
			workspacePath: createUniquePath("./tmp"),
			nodeName:      "example",
			want:          false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			metadataFilePath := createUniquePath("./.notewolfy")
			config := &structure.Config{
				MetadataFilePath: metadataFilePath,
			}
			defer CleanUpFile(metadataFilePath)

			mmf, err := structure.NewMetadataNoteWolfyFileHandle(config)
			assert.NoError(t, err)

			// We need to prepare the workspace
			workspacePath, err := utility.ExpandRelativePaths(tc.workspacePath)
			assert.NoError(t, err)
			err = os.Mkdir(workspacePath, os.ModePerm)
			assert.NoError(t, err)
			workspaceNode := &structure.Node{
				Name: "Workspace",
				Path: workspacePath,
			}
			mmf.Workspaces = append(mmf.Workspaces, workspaceNode)
			mmf.ActiveWorkspace = workspaceNode.Name
			err = mmf.Save()
			assert.NoError(t, err)

			// We need to create a node
			nodePath := filepath.Join(tc.workspacePath, tc.nodeName)
			err = os.Mkdir(nodePath, os.ModePerm)
			assert.NoError(t, err)
			node := &structure.Node{
				Name: tc.nodeName,
				Path: nodePath,
			}
			mmf.ActiveNode = node.Name
			mmf.Workspaces[0].Children = append(mmf.Workspaces[0].Children, node)
			err = mmf.Save()
			assert.NoError(t, err)

			defer os.Remove(workspacePath)
			defer os.Remove(nodePath)

			statement := tc.statement
			commands.MatchStatementToCommand(mmf, statement)
			if tc.want {
				assert.Equal(t, workspaceNode.Name, mmf.ActiveNode)

				return
			}
			assert.Equal(t, node.Name, mmf.ActiveNode)
		})
	}
}

func TestMatchStatementToOpen(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		statement     string
		workspacePath string
		workspaceName string
		nodeName      string
		want          bool
	}{
		"simple open command": {
			statement:     "open test",
			workspacePath: createUniquePath("./tmp"),
			workspaceName: "test",
			nodeName:      "example",
			want:          true,
		},
		"error open command": {
			statement:     "opens test",
			workspacePath: createUniquePath("./tmp"),
			workspaceName: "test",
			nodeName:      "example",
			want:          false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			metadataFilePath := createUniquePath("./.notewolfy")
			config := &structure.Config{
				MetadataFilePath: metadataFilePath,
			}
			defer CleanUpFile(metadataFilePath)

			mmf, err := structure.NewMetadataNoteWolfyFileHandle(config)
			assert.NoError(t, err)

			// We need to prepare the workspace
			workspacePath, err := utility.ExpandRelativePaths(tc.workspacePath)
			assert.NoError(t, err)
			err = os.Mkdir(workspacePath, os.ModePerm)
			assert.NoError(t, err)
			workspaceNode := &structure.Node{
				Name: tc.workspaceName,
				Path: workspacePath,
			}
			mmf.Workspaces = append(mmf.Workspaces, workspaceNode)
			mmf.ActiveWorkspace = workspaceNode.Name
			err = mmf.Save()
			assert.NoError(t, err)

			// We need to create a node
			nodePath := filepath.Join(tc.workspacePath, tc.nodeName)
			err = os.Mkdir(nodePath, os.ModePerm)
			assert.NoError(t, err)
			node := &structure.Node{
				Name: tc.nodeName,
				Path: nodePath,
			}
			mmf.ActiveNode = node.Name
			mmf.Workspaces[0].Children = append(mmf.Workspaces[0].Children, node)
			err = mmf.Save()
			assert.NoError(t, err)

			defer os.Remove(workspacePath)
			defer os.Remove(nodePath)

			statement := tc.statement
			commands.MatchStatementToCommand(mmf, statement)
			if tc.want {
				assert.Equal(t, workspaceNode.Name, mmf.ActiveWorkspace)
				assert.Equal(t, workspaceNode.Name, mmf.ActiveNode)

				return
			}
			assert.Equal(t, node.Name, mmf.ActiveNode)
		})
	}
}

func TestMatchStatementToEdit(t *testing.T) {
	tests := map[string]struct {
		statement     string
		workspacePath string
		want          bool
	}{
		"simple edit command": {
			statement:     "edit example",
			workspacePath: createUniquePath("./tmp"),
			want:          true,
		},
		"error edit command": {
			statement:     "edits example",
			workspacePath: createUniquePath("./tmp"),
			want:          false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			metadataFilePath := createUniquePath("./.notewolfy")
			config := &structure.Config{
				MetadataFilePath: metadataFilePath,
			}
			defer CleanUpFile(metadataFilePath)

			mmf, err := structure.NewMetadataNoteWolfyFileHandle(config)
			assert.NoError(t, err)

			// We need to prepare the workspace
			workspacePath, err := utility.ExpandRelativePaths(tc.workspacePath)
			assert.NoError(t, err)
			err = os.Mkdir(workspacePath, os.ModePerm)
			assert.NoError(t, err)
			workspaceNode := &structure.Node{
				Name: "workspace",
				Path: workspacePath,
			}
			mmf.ActiveNode = workspaceNode.Name
			mmf.Workspaces = append(mmf.Workspaces, workspaceNode)
			mmf.ActiveWorkspace = workspaceNode.Name
			err = mmf.Save()
			assert.NoError(t, err)

			// We need to create a markdown file
			markdownFileName := "example.md"
			markdownFilePath := filepath.Join(workspacePath, markdownFileName)
			markdown := &structure.Markdown{
				Filename: markdownFileName,
			}
			file, err := os.Create(markdownFilePath)
			assert.NoError(t, err)
			mmf.Workspaces[0].Markdowns = append(mmf.Workspaces[0].Markdowns, markdown)
			err = mmf.Save()
			assert.NoError(t, err)

			defer os.Remove(workspacePath)
			defer func() {
				file.Close()
				os.Remove(file.Name())
			}()

			tempFile, err := os.CreateTemp("", "tempStdin")
			assert.NoError(t, err)
			defer os.Remove(tempFile.Name())

			_, err = tempFile.WriteString("iHelloWorld!\nThis is a Test!\x1b:wq\n")
			assert.NoError(t, err)
			err = tempFile.Close()
			assert.NoError(t, err)

			tempFile, err = os.Open(tempFile.Name())
			assert.NoError(t, err)
			defer tempFile.Close()

			oldStdin := os.Stdin
			defer func() { os.Stdin = oldStdin }()
			os.Stdin = tempFile

			statement := tc.statement
			_, err = captureStdOutput(func() {
				commands.MatchStatementToCommand(mmf, statement)
			})
			assert.NoError(t, err)
			if tc.want {
				content, err := os.ReadFile(markdownFilePath)
				assert.NoError(t, err)
				actContentString := string(content)
				expContentString := "HelloWorld!\nThis is a Test!\n"
				assert.Equal(t, expContentString, actContentString)

				return
			}
		})
	}
}

func TestMatchStatementToVersion(t *testing.T) {
	tests := map[string]struct {
		statement string
		want      bool
	}{
		"simple version command": {
			statement: "version",
			want:      true,
		},
		"error version command": {
			statement: "versions",
			want:      false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			metadataFilePath := createUniquePath("./.notewolfy")
			config := &structure.Config{
				MetadataFilePath: metadataFilePath,
			}
			defer CleanUpFile(metadataFilePath)

			mmf, err := structure.NewMetadataNoteWolfyFileHandle(config)
			assert.NoError(t, err)
			actOutput, err := captureStdOutput(func() {
				commands.MatchStatementToCommand(mmf, tc.statement)
			})
			assert.NoError(t, err)
			if tc.want {
				expContentString := "\n\rnotewolfy version v0.2.0 at your disposal!"
				assert.Equal(t, expContentString, actOutput)

				return
			}
		})
	}
}

func TestMatchStatementToList(t *testing.T) {
	tests := map[string]struct {
		statement     string
		workspacePath string
		workspaceName string
		want          bool
	}{
		"simple list command": {
			statement:     "ls",
			workspacePath: createUniquePath("./tmp"),
			workspaceName: "test",
			want:          true,
		},
		"error list command": {
			statement:     "lss",
			workspacePath: createUniquePath("./tmp"),
			workspaceName: "test",
			want:          false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			metadataFilePath := createUniquePath("./.notewolfy")
			config := &structure.Config{
				MetadataFilePath: metadataFilePath,
			}
			defer CleanUpFile(metadataFilePath)

			mmf, err := structure.NewMetadataNoteWolfyFileHandle(config)
			assert.NoError(t, err)

			// We need to prepare the workspace
			workspacePath, err := utility.ExpandRelativePaths(tc.workspacePath)
			assert.NoError(t, err)
			err = os.Mkdir(workspacePath, os.ModePerm)
			assert.NoError(t, err)
			workspaceNode := &structure.Node{
				Name: tc.workspaceName,
				Path: workspacePath,
			}
			mmf.Workspaces = append(mmf.Workspaces, workspaceNode)
			mmf.ActiveNode = workspaceNode.Name
			mmf.ActiveWorkspace = workspaceNode.Name
			err = mmf.Save()
			assert.NoError(t, err)

			// We need to create nodes and markdown files that ls can display
			nodePathA := filepath.Join(tc.workspacePath, "A")
			nodeA := &structure.Node{
				Name: "A",
				Path: nodePathA,
			}
			nodePathB := filepath.Join(tc.workspacePath, "B")
			nodeB := &structure.Node{
				Name: "B",
				Path: nodePathB,
			}
			markdown := &structure.Markdown{
				Filename: "example.md",
			}
			mmf.Workspaces[0].Children = append(mmf.Workspaces[0].Children, nodeA)
			mmf.Workspaces[0].Children = append(mmf.Workspaces[0].Children, nodeB)
			mmf.Workspaces[0].Markdowns = append(mmf.Workspaces[0].Markdowns, markdown)
			err = mmf.Save()
			assert.NoError(t, err)

			defer os.Remove(workspacePath)

			statement := tc.statement
			actOutput, err := captureStdOutput(func() {
				commands.MatchStatementToCommand(mmf, statement)
			})
			assert.NoError(t, err)
			if tc.want {
				expOutput := "\r\nYou are on node:  test\n\rChild nodes:\n\r A\n\r B\n\rMarkdown files:\n\r example.md\n"
				assert.Equal(t, expOutput, actOutput)

				return
			}
			assert.Empty(t, actOutput)
		})
	}
}

func TestMatchStatementToListWorkspaces(t *testing.T) {
	// Note: These tests only work with ./ as a base path
	tests := map[string]struct {
		statement      string
		workspacePaths []string
		workspaceName  string
		want           bool
	}{
		"simple ls ws command": {
			statement:      "ls ws",
			workspacePaths: []string{createUniquePath("./tmpA"), createUniquePath("./tmpB")},
			workspaceName:  "test",
			want:           true,
		},
		"error ls ws command": {
			statement:      "ls s ws",
			workspacePaths: []string{createUniquePath("./tmpA"), createUniquePath("./tmpB")},
			workspaceName:  "test",
			want:           false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			metadataFilePath := createUniquePath("./.notewolfy")
			config := &structure.Config{
				MetadataFilePath: metadataFilePath,
			}
			defer CleanUpFile(metadataFilePath)

			mmf, err := structure.NewMetadataNoteWolfyFileHandle(config)
			assert.NoError(t, err)

			// We need to prepare the workspaces
			workspacePathA, err := utility.ExpandRelativePaths(tc.workspacePaths[0])
			assert.NoError(t, err)
			workspaceNodeA := &structure.Node{
				Name: tc.workspaceName + "A",
				Path: workspacePathA,
			}
			workspacePathB, err := utility.ExpandRelativePaths(tc.workspacePaths[1])
			assert.NoError(t, err)
			workspaceNodeB := &structure.Node{
				Name: tc.workspaceName + "B",
				Path: workspacePathB,
			}
			mmf.Workspaces = append(mmf.Workspaces, workspaceNodeA)
			mmf.Workspaces = append(mmf.Workspaces, workspaceNodeB)
			mmf.ActiveNode = workspaceNodeA.Name
			mmf.ActiveWorkspace = workspaceNodeA.Name
			err = mmf.Save()
			assert.NoError(t, err)

			defer os.Remove(workspacePathA)
			defer os.Remove(workspacePathB)

			statement := tc.statement
			actOutput, err := captureStdOutput(func() {
				commands.MatchStatementToCommand(mmf, statement)
			})
			assert.NoError(t, err)
			if tc.want {
				basePath, err := utility.ExpandRelativePaths("./")
				assert.NoError(t, err)
				workspaceNameA := tc.workspacePaths[0][2:]
				workspaceNameB := tc.workspacePaths[1][2:]
				expOutput := fmt.Sprintf("\r\nWorkspace Name                                                                                              Workspace Path                                                                                          \r\n-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------\n\rtestA                                                                                                       %[1]s/%[2]s\n\rtestB                                                                                                       %[1]s/%[3]s\n", basePath, workspaceNameA, workspaceNameB)
				assert.Equal(t, expOutput, actOutput)

				return
			}
			expOutput := "\r\nYou are on node:  testA\n\rChild nodes:\n\rMarkdown files:\n"
			assert.Equal(t, expOutput, actOutput)
		})
	}
}

func TestMatchStatementToHelp(t *testing.T) {
	tests := map[string]struct {
		statement string
		expOutput string
	}{
		"simple help command (1)": {
			statement: "help ls",
			expOutput: "\n\rCommand: ls\n\rDescription: ls can be used to list information about the node that you are on, e.g. active node, markdown files on that node, etc.\n\rExample Usage: ls",
		},
		"simple help command (2)": {
			statement: "help create workspace",
			expOutput: "\n\rCommand: create workspace <workspaceName> <workspacePath>\n\rDescription: create workspace will create a new workspace for you under the specified name and path that you can choose.\n\rExample Usage: create workspace example /path/to/example",
		},
		"error help command": {
			statement: "help something",
			expOutput: "\n\rYou need to specify a valid command, here is a list of possible commands:\n\r- ls\n\r- ls ws\n\r- create workspace\n\r- delete workspace\n\r- create node\n\r- delete node\n\r- create md\n\r- delete md\n\r- edit\n\r- goto\n\r- goback\n\r- open\n\r- version",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			metadataFilePath := createUniquePath("./.notewolfy")
			config := &structure.Config{
				MetadataFilePath: metadataFilePath,
			}
			defer CleanUpFile(metadataFilePath)

			mmf, err := structure.NewMetadataNoteWolfyFileHandle(config)
			assert.NoError(t, err)

			actOutput, err := captureStdOutput(func() {
				commands.MatchStatementToCommand(mmf, tc.statement)
			})
			assert.NoError(t, err)
			assert.Equal(t, tc.expOutput, actOutput)
		})
	}
}

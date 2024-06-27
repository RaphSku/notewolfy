package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/RaphSku/notewolfy/internal/structure"
)

type CreateMarkdownStrategy struct {
	statement string
	mmf       *structure.MetadataNoteWolfyFileHandle
}

func (cms *CreateMarkdownStrategy) Run() error {
	markdownNameRegex := regexp.MustCompile("create md (?P<name>[\\w]+)")
	matches := markdownNameRegex.FindStringSubmatch(cms.statement)
	names := markdownNameRegex.SubexpNames()
	namedGroups := make(map[string]string)
	for i, name := range names {
		if i != 0 && name != "" {
			namedGroups[name] = matches[i]
		}
	}
	markdownName := namedGroups["name"]
	markdownNameWithFExt := strings.Join([]string{markdownName, ".md"}, "")
	activeNodeName := cms.mmf.ActiveNode
	activeNode := cms.mmf.FindNode(activeNodeName)
	pathToMarkdown := filepath.Join(activeNode.Path, markdownNameWithFExt)

	_, err := os.Stat(pathToMarkdown)
	if os.IsNotExist(err) {
		markdown := &structure.Markdown{
			Filename: markdownNameWithFExt,
		}
		cms.mmf.AddMarkdown(markdown)
		cms.mmf.Save()

		file, err := os.Create(pathToMarkdown)
		if err != nil {
			return err
		}
		defer file.Close()

		return nil
	}

	fmt.Println("\r\nMarkdown file already exists!")

	return nil
}

type DeleteMDStrategy struct {
	statement string
	mmf       *structure.MetadataNoteWolfyFileHandle
}

func (dms *DeleteMDStrategy) Run() error {
	markdownNameRegex := regexp.MustCompile("delete md (?P<name>[\\w]+)")
	matches := markdownNameRegex.FindStringSubmatch(dms.statement)
	names := markdownNameRegex.SubexpNames()
	namedGroups := make(map[string]string)
	for i, name := range names {
		if i != 0 && name != "" {
			namedGroups[name] = matches[i]
		}
	}
	markdownName := namedGroups["name"]

	activeNodeName := dms.mmf.ActiveNode
	activeNode := dms.mmf.FindNode(activeNodeName)
	err := os.Remove(filepath.Join(activeNode.Path, strings.Join([]string{markdownName, ".md"}, "")))
	if err != nil {
		return err
	}

	err = dms.mmf.DeleteMarkdown(markdownName)
	if err != nil {
		return err
	}
	dms.mmf.Save()

	return nil
}

type EditStrategy struct {
	statement string
	mmf       *structure.MetadataNoteWolfyFileHandle
}

func (es *EditStrategy) Run() error {
	markdownNameRegex := regexp.MustCompile("edit (?P<name>[\\w]+)")
	matches := markdownNameRegex.FindStringSubmatch(es.statement)
	names := markdownNameRegex.SubexpNames()
	namedGroups := make(map[string]string)
	for i, name := range names {
		if i != 0 && name != "" {
			namedGroups[name] = matches[i]
		}
	}
	markdownName := namedGroups["name"]
	activeNodeName := es.mmf.ActiveNode
	activeNode := es.mmf.FindNode(activeNodeName)
	for _, markdown := range activeNode.Markdowns {
		if markdown.Filename[:len(markdown.Filename)-3] == markdownName {
			markdownFile := filepath.Join(activeNode.Path, markdown.Filename)

			cmd := exec.Command("vi", markdownFile)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err := cmd.Run()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

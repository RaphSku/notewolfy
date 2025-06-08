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
	nameCaptureGroupName := "name"
	markdownNamePattern := "[\\w]+"
	pattern := fmt.Sprintf("create md (?P<%s>%s)", nameCaptureGroupName, markdownNamePattern)
	markdownNameRegex := regexp.MustCompile(pattern)
	matches := markdownNameRegex.FindStringSubmatch(cms.statement)
	if len(matches) != 2 {
		return fmt.Errorf("\n\rPlease check whether the markdown name matches the regex %s!", markdownNamePattern)
	}
	names := markdownNameRegex.SubexpNames()
	var markdownName string
	for i, name := range names[1:] {
		if name == nameCaptureGroupName {
			markdownName = matches[i+1]
		}
	}
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

	return fmt.Errorf("\r\nMarkdown file already exists!")
}

type DeleteMDStrategy struct {
	statement string
	mmf       *structure.MetadataNoteWolfyFileHandle
}

func (dms *DeleteMDStrategy) Run() error {
	nameCaptureGroupName := "name"
	markdownNamePattern := "[\\w]+"
	pattern := fmt.Sprintf("delete md (?P<%s>%s)", nameCaptureGroupName, markdownNamePattern)
	markdownNameRegex := regexp.MustCompile(pattern)
	matches := markdownNameRegex.FindStringSubmatch(dms.statement)
	if len(matches) != 2 {
		return fmt.Errorf("\n\rPlease check whether the markdown name matches the regex %s!", markdownNamePattern)
	}
	names := markdownNameRegex.SubexpNames()
	var markdownName string
	for i, name := range names[1:] {
		if name == nameCaptureGroupName {
			markdownName = matches[i+1]
		}
	}

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
	nameCaptureGroupName := "name"
	markdownNamePattern := "[\\w]+"
	pattern := fmt.Sprintf("edit (?P<%s>%s)", nameCaptureGroupName, markdownNamePattern)
	markdownNameRegex := regexp.MustCompile(pattern)
	matches := markdownNameRegex.FindStringSubmatch(es.statement)
	if len(matches) != 2 {
		return fmt.Errorf("\n\rPlease check whether the markdown name matches the regex %s!", markdownNamePattern)
	}
	names := markdownNameRegex.SubexpNames()
	var markdownName string
	for i, name := range names[1:] {
		if name == nameCaptureGroupName {
			markdownName = matches[i+1]
		}
	}
	activeNodeName := es.mmf.ActiveNode
	activeNode := es.mmf.FindNode(activeNodeName)
	for _, markdown := range activeNode.Markdowns {
		if markdown.Filename[:len(markdown.Filename)-3] == markdownName {
			markdownFile := filepath.Join(activeNode.Path, markdown.Filename)

			cmd := exec.Command("vim", markdownFile)
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

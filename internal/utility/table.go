package utility

import (
	"errors"
	"fmt"
	"strings"
)

type Table struct {
	width   int
	headers []string
	entries [][]string
}

func NewTable(width int) *Table {
	return &Table{
		width: width,
	}
}

func (t *Table) SetHeaders(headers []string) {
	t.headers = append(t.headers, headers...)
}

func (t *Table) Append(entriesInRow []string) error {
	if len(entriesInRow) != len(t.headers) {
		return errors.New("could not append entries because the length does not match")
	}
	t.entries = append(t.entries, entriesInRow)

	return nil
}

func (t *Table) GenerateTable() error {
	if len(t.headers) == 0 {
		return errors.New("please set headers in order to generate the table")
	}

	for index, header := range t.headers {
		centeredHeader := t.centerString(header, t.width)
		fmt.Print(centeredHeader)
		if index != len(t.headers)-1 {
			fmt.Print("|")
		}
	}
	fmt.Print("\r\n")
	divisionLine := t.divisionLine()
	fmt.Print(divisionLine)
	fmt.Print("\r\n")
	for _, entriesInRow := range t.entries {
		for j, entry := range entriesInRow {
			centeredEntry := t.centerString(entry, t.width)
			fmt.Print(centeredEntry)
			if j != len(t.headers)-1 {
				fmt.Print("|")
			}
		}
		fmt.Print("\r\n")
	}
	fmt.Print("\r\n")

	return nil
}

func (t *Table) centerString(target string, width int) string {
	padding := width - len(target)
	if padding <= 0 {
		return target
	}
	leftPadding := padding / 2
	rightPadding := padding - leftPadding

	return strings.Repeat(" ", leftPadding) + target + strings.Repeat(" ", rightPadding)
}

func (t *Table) divisionLine() string {
	totalWidth := len(t.headers)*t.width + len(t.headers) - 1

	return strings.Repeat("-", totalWidth)
}

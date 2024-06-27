//go:build unit_test

package utility_test

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/RaphSku/notewolfy/internal/utility"
	"github.com/stretchr/testify/assert"
)

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

func TestGenerateTable(t *testing.T) {
	expHeaders := []string{"testA", "testB", "testC"}
	expEntries := [][]string{
		{"A", "B", "C"},
		{"hello", "world", "everybody"},
	}

	table := utility.NewTable(20)
	table.SetHeaders(expHeaders)
	for _, entryInRow := range expEntries {
		err := table.Append(entryInRow)
		assert.NoError(t, err)
	}

	actTable, err := captureStdOutput(func() {
		err := table.GenerateTable()
		assert.NoError(t, err)
	})
	assert.NoError(t, err)
	expTable := "       testA        |       testB        |       testC        \r\n--------------------------------------------------------------\r\n         A          |         B          |         C          \r\n       hello        |       world        |     everybody      \r\n\r\n"
	assert.Equal(t, expTable, actTable)
}

func TestEntryLengthNotMatching(t *testing.T) {
	headers := []string{"testA", "testB", "testC"}
	faultyEntry := []string{"A", "B"}

	table := utility.NewTable(20)
	table.SetHeaders(headers)

	err := table.Append(faultyEntry)
	expError := errors.New("could not append entries because the length does not match")
	if assert.Error(t, err) {
		assert.Equal(t, expError, err)
	}
}

func TestGenerateTableWithNoHeaders(t *testing.T) {
	table := utility.NewTable(20)
	err := table.GenerateTable()
	expError := errors.New("please set headers in order to generate the table")
	if assert.Error(t, err) {
		assert.Equal(t, expError, err)
	}
}

func TestCenterStringWithTooNarrowPadding(t *testing.T) {
	headers := []string{"testA", "testB", "testC"}
	entry := []string{"And here", "we go again", "too bad"}
	table := utility.NewTable(5)
	table.SetHeaders(headers)
	table.Append(entry)
	actOutput, err := captureStdOutput(func() {
		err := table.GenerateTable()
		assert.NoError(t, err)
	})
	assert.NoError(t, err)
	expOutput := "testA|testB|testC\r\n-----------------\r\nAnd here|we go again|too bad\r\n\r\n"
	assert.Equal(t, expOutput, actOutput)
}

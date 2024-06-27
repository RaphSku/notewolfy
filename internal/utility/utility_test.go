//go:build unit_test

package utility_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/RaphSku/notewolfy/internal/utility"
	"github.com/stretchr/testify/assert"
)

func TestGetHomeDir(t *testing.T) {
	t.Parallel()

	expHomeDir, err := os.UserHomeDir()
	assert.NoError(t, err)

	actHomeDir, err := utility.GetHomeDir()
	assert.NoError(t, err)

	assert.Equal(t, expHomeDir, actHomeDir)
}

func TestExpandRelativePaths(t *testing.T) {
	t.Parallel()

	expHomeDir, err := os.UserHomeDir()
	assert.NoError(t, err)

	expTildePath := filepath.Join(expHomeDir, "/tmp")
	expDotPath, err := filepath.Abs("./tmp")
	assert.NoError(t, err)
	expPrevPath, err := filepath.Abs("../tmp")
	assert.NoError(t, err)

	testCases := []struct {
		name      string
		givenPath string
		expPath   string
	}{
		{name: "Check tilde relative path", givenPath: "~/tmp", expPath: expTildePath},
		{name: "Check dot relative path", givenPath: "./tmp", expPath: expDotPath},
		{name: "Check previous directory relative path", givenPath: "../tmp", expPath: expPrevPath},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actPath, err := utility.ExpandRelativePaths(tc.givenPath)
			assert.NoError(t, err)
			assert.Equal(t, tc.expPath, actPath)
		})
	}
}

func TestDoesChildPathMatchesParentPath(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		givenPaths []string
		expResult  bool
		expError   bool
	}{
		{name: "Check paths matching with tilde", givenPaths: []string{"~/tmp", "~/tmp/A"}, expResult: true, expError: false},
		{name: "Check paths matching with dot", givenPaths: []string{"./tmp", "./tmp/A"}, expResult: true, expError: false},
		{name: "Check paths matching", givenPaths: []string{"/A/B", "/A/B/C"}, expResult: true, expError: false},
		{name: "Check paths do not match", givenPaths: []string{"/A", "/B"}, expResult: false, expError: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ok, err := utility.DoesChildPathMatchesParentPath(tc.givenPaths[0], tc.givenPaths[1])
			if tc.expError {
				if assert.Error(t, err) {
					expError := errors.New("child's parent path and parentPath do not match")
					assert.Equal(t, expError, err)
				}
				return
			}
			assert.Equal(t, tc.expResult, ok)
			assert.NoError(t, err)
		})
	}
}

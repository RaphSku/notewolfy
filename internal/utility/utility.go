package utility

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
)

func GetHomeDir() (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return userHomeDir, nil
}

func ExpandRelativePaths(path string) (string, error) {
	expandedPath := path
	if path[0] == '~' {
		homeDir, err := GetHomeDir()
		if err != nil {
			return "", err
		}
		expandedPath = filepath.Join(homeDir, path[1:])
	} else if path[:2] == ".." || path[0] == '.' {
		var err error
		expandedPath, err = filepath.Abs(path)
		if err != nil {
			return "", err
		}
	}

	return expandedPath, nil
}

func DoesChildPathMatchesParentPath(parentPath string, childPath string) (bool, error) {
	re, err := regexp.Compile("/[\\w.]+")
	if err != nil {
		return false, err
	}

	matches := re.FindAllString(childPath, -1)
	childDirectoryLength := len(matches[len(matches)-1])
	childsParentPath := childPath[:len(childPath)-childDirectoryLength]
	if childsParentPath != parentPath {
		return false, errors.New("child's parent path and parentPath do not match")
	}

	return true, nil
}

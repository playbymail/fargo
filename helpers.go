// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package fargo

import (
	"os"
	"path/filepath"
)

func AbsPath(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	} else if sb, err := os.Stat(absPath); err != nil {
		return "", err
	} else if !sb.IsDir() {
		return "", ErrNotADirectory
	}
	return absPath, nil
}

func IsDir(path string) (bool, error) {
	sb, err := os.Stat(path)
	if err != nil {
		return false, err
	} else if !sb.IsDir() {
		return false, nil
	}
	return true, nil
}

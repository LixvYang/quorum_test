package utils

import (
	"os"

	logging "github.com/ipfs/go-log/v2"
)

var logger = logging.Logger("utils")

// FileExist check if file is exist
func FileExist(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// DirExist check if file is exist
func DirExist(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return info.IsDir()
}

// EnsureDir make sure `dir` exist ,or create it
func EnsureDir(dir string) error {
	if !DirExist(dir) {
		logger.Infof("try to create directory: %s", dir)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			logger.Infof("make directory %s failed: %s", dir, err)
			return err
		}
	}
	return nil
}

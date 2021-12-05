package utils

import (
	"math/rand"
	"os"
	"io"

	logging "github.com/ipfs/go-log/v2"
	maddr "github.com/multiformats/go-multiaddr"
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

func GetRandomStr(n int) string {
	const letters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const lettersLength = int64(len(letters))

	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Int63()%lettersLength]
	}
	return string(b)
}

func StringsToAddrs(addrStrings []string) (maddrs []maddr.Multiaddr, err error) {
	for _, addrString := range addrStrings {
		addr, err := maddr.NewMultiaddr(addrString)
		if err != nil {
			return maddrs, err
		}
		maddrs = append(maddrs, addr)
	}
	return
}

// Is DirEmpty check if dir is empty
func IsDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err //
}
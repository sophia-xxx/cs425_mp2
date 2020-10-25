package file_manager

import (
	"io/ioutil"
	"os"

	"cs425_mp2/config"
)

func GetSDFSFilePath(filename string) string {
	return config.SDFS_DIR + filename
}

func WhetherFileExist(filepath string) bool {
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// when failed node rejoin the system, it will remove all sdfs files
func RemoveAllSDFSFile() {
	os.RemoveAll(config.SDFS_DIR)
	os.MkdirAll(config.SDFS_DIR, config.PERM_MODE)
}

func RemoveSDFSFile(filename string) {
	os.Remove(config.SDFS_DIR + filename)
}

func GetLocalSDFSFileList() []string {
	result := make([]string, 0)
	files, _ := ioutil.ReadDir(config.SDFS_DIR)
	for _, file := range files {
		result = append(result, file.Name())
	}
	return result
}

package file_manager

import (
	config2 "cs425_mp2/config"
	"cs425_mp2_remastered/config"
	"io/ioutil"
	"os"
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
	os.RemoveAll(config2.SDFS_DIR)
	os.MkdirAll(config2.SDFS_DIR, config2.PERM_MODE)
}

func RemoveSDFSFile(filename string) {
	os.Remove(config2.SDFS_DIR + filename)
}

func GetAllSDFSFiles() []os.FileInfo {
	files, _ := ioutil.ReadDir(config.SDFS_DIR)
	return files
}

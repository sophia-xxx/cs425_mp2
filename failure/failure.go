package failure

import (
	"cs425_mp2/config"
	"os"
)

// when failed node rejoin the system, it will remove all sdfs files
func RemoveAllFile() {
	os.RemoveAll(config.SDFS_DIR)
	os.MkdirAll(config.SDFS_DIR, config.PERM_MODE)
}

// master election
func election() {

}

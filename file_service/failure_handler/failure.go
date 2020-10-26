package failure_handler

import (
	"cs425_mp2/file_service/file_manager"
	"cs425_mp2/file_service/file_record"
	"cs425_mp2/file_service/networking"
	"cs425_mp2/member_service"
	"cs425_mp2/util/logger"
)

// sendLocalSDFSFileInfo To new Master
func HandleMasterFailure() {
	select {
	case <-member_service.MasterChanged:
		logger.PrintInfo("File service noticed that master changed. Handling Master failure...")
		if member_service.IsMaster() {
			file_record.NewMasterInit()
		} else {
			uploadSDFSFileListToMaster()
		}
	default:
		return
	}
}

func uploadSDFSFileListToMaster() {
	logger.PrintInfo("Restoring file record to new master...")
	networking.RestoreFileListToMaster(file_manager.GetLocalSDFSFileList(), member_service.GetMasterIP())
}
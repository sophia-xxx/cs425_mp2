package failure_handler

import (
	"cs425_mp2/file_service/file_manager"
	"cs425_mp2/file_service/file_record"
	"cs425_mp2/file_service/networking"
	"cs425_mp2/member_service"
)

// sendLocalSDFSFileInfo To new Master
func HandleMasterFailure() {
	select {
	case <-member_service.MasterChanged:
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
	networking.RestoreFileListToMaster(file_manager.GetLocalSDFSFileList(), member_service.GetMasterIP())
}
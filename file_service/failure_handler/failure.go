package failure_handler

import "cs425_mp2/member_service"

// sendLocalSDFSFileInfo To new Master
func HandleMasterFailure() {
	select {
	case <-member_service.MasterChanged:
		uploadSDFSFileToMaster()
	default:
		return
	}
}

func uploadSDFSFileToMaster() {

}
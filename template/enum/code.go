package enum

type StatusCode int

const (
	Success           StatusCode = 200001
	InvalidParams     StatusCode = 400001
	FailedAdminCtx    StatusCode = 400002
	InvalidFileFormat StatusCode = 410001
	FailedUploadFile  StatusCode = 410002
	FailedImportFile  StatusCode = 410003
	FailedExportFile  StatusCode = 410004
	FailedGetData     StatusCode = 420001
	FailedCreateData  StatusCode = 420002
	FailedUpdateData  StatusCode = 420003
	FailedDeleteData  StatusCode = 420004
	FailedSearchData  StatusCode = 420005
)

var (
	statusCodeMsg = map[StatusCode]string{
		Success:           "操作成功",
		InvalidParams:     "参数错误",
		FailedAdminCtx:    "创建管理员上下文失败",
		InvalidFileFormat: "文件格式错误",
		FailedUploadFile:  "上传文件失败",
		FailedImportFile:  "导入文件失败",
		FailedExportFile:  "导出文件失败",
		FailedGetData:     "获取数据失败",
		FailedCreateData:  "创建数据失败",
		FailedUpdateData:  "更新数据失败",
		FailedDeleteData:  "删除数据失败",
		FailedSearchData:  "搜索数据失败",
	}
)

func (sc StatusCode) String() string {
	return statusCodeMsg[sc]
}

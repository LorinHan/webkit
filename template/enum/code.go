package enum

type StatusCode = int

const (
	Success StatusCode = 200001 // 成功

	InvalidParams  StatusCode = 400001 // 参数错误
	FailedAdminCtx StatusCode = 400002 // 创建管理员上下文失败

	InvalidFileFormat StatusCode = 410001 // 文件格式错误
	FailedUploadFile  StatusCode = 410002 // 上传文件失败
	FailedImportFile  StatusCode = 410003 // 导入文件失败
	FailedExportFile  StatusCode = 410004 // 导出文件失败

	FailedGetData    StatusCode = 420001 // 获取数据失败
	FailedCreateData StatusCode = 420002 // 创建数据失败
	FailedUpdateData StatusCode = 420003 // 更新数据失败
	FailedDeleteData StatusCode = 420004 // 删除数据失败
	FailedSearchData StatusCode = 420005 // 搜索数据失败
)

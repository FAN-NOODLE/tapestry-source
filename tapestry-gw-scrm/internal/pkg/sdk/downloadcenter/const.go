package downloadcenter

const (
	UrlCreateExport       = "/v1/create-export"
	UrlExportList         = "/v1/export-list"
	UrlRetry              = "/v1/retry"
	UrlProcess            = "/v1/process"
	UrlDownload           = "/v1/download"
	UrlDelete             = "/v1/delete"
	UrlDetail             = "/v1/detail"
	UrlLatestExportDetail = "/v1/latest-export"
	UrlAutoDetail         = "/v1/auto-detail"
)

const (
	FileTypeExcel = 1
	FileTypeCsv   = 2
	FileTypeZip   = 3
)
const (
	// RestfulPaginationTypePage Restful分页类型: 页码
	RestfulPaginationTypePage = 1
	// RestfulPaginationTypeID Restful分页类型: ID
	RestfulPaginationTypeID = 2
)

const (
	SourceTypeData    = 1
	SourceTypeSql     = 2
	SourceTypeRestful = 3
)

// 标签常量
const (
	// ReserveOrderExport 邀约单列表导出
	ReserveOrderExport = "reserve-order-export"
	VmImageSave        = "vm-image-save"
	VmListExportSave   = "vm-list-export"
)

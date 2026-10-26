package downloadcenter

import "gw-scrm/internal/pkg/sdk"

type baseResponse struct {
	sdk.ResponseBase
	Message string `json:"message"`
}

// CreateData 数据类型
type CreateData struct {
	SourceData [][]interface{} `json:"source_data"`
}

// CreateSql sql类型
type CreateSql struct {
	SourceSqlDatabase string `json:"source_sql_database"`
	SourceSqlCount    uint8  `json:"source_sql_count"`
	SourceSql         string `json:"source_sql"`
}

// CreateRestful Restful类型
type CreateRestful struct {
	SourceRestfulPath    string                 `json:"source_restful_path"`
	SourceRestfulMethod  string                 `json:"source_restful_method"`
	SourceRestfulHeaders map[string]interface{} `json:"source_restful_headers"`
	SourceRestfulParams  map[string]interface{} `json:"source_restful_params"`

	SourceSuffix                  []string `json:"source_suffix"`
	SourceRestfulCountField       string   `json:"source_restful_count_field"`
	SourceRestfulDataField        string   `json:"source_restful_data_field"`
	SourceRestfulRequestPageField string   `json:"source_restful_request_page_field"`
	SourceRestfulRequestSizeField string   `json:"source_restful_request_size_field"`
	SourceExcelCellMerge          struct {
		Type    uint8    `json:"type"`
		Rule    string   `json:"rule"`
		Title   []string `json:"title"`
		Premise []string `json:"premise"`
	} `json:"source_excel_cell_merge"`
	SourceRestfulPaginationType int               `json:"source_restful_pagination_type"` // 分页类型：1:页码 2:ID
	SourceRestfulMinId          int64             `json:"source_restful_min_id"`
	SourceRestfulMaxId          int64             `json:"source_restful_max_id"`
	SourceRestfulTotal          int               `json:"source_restful_total"`
	SourceRestfulExt            map[string]string `json:"source_restful_ext"`
}

type CreatePack struct {
}

type CreateExportRequest struct {
	FileName           string            `json:"filename"`
	Type               uint8             `json:"type"`
	StoreId            int64             `json:"store_id"`
	EnterpriseId       int64             `json:"enterprise_id"`
	OperationId        int64             `json:"operation_id"`
	OperationName      string            `json:"operation_name"`
	SourceType         uint8             `json:"source_type"`
	SourceTitleAlias   map[string]string `json:"source_title_alias"`
	SourceTitle        []string          `json:"source_title"`
	SourceTitleValue   []string          `json:"source_title_value" comment:"字段的值不需要转成string,在这填写"`
	SourceHeaderPrefix [][]string        `json:"source_header_prefix"`
	Label              string            `json:"label"`

	CreateData
	CreateSql
	CreateRestful
}

type CreateExportResponse struct {
	baseResponse

	Data struct {
		FileId int64 `json:"file_id"`
	} `json:"data"`
}

type ExportListRequest struct {
	Page          int     `json:"page"`
	PageSize      int     `json:"page_size"`
	Status        uint8   `json:"status"`
	Filename      string  `json:"filename"`
	EnterpriseId  int64   `json:"enterprise_id"`
	OperationIds  []int64 `json:"operation_ids"`
	OperationId   int64   `json:"operation_id"`
	OperationName string  `json:"operation_name"`
	CreatedAt     struct {
		Min int64 `json:"min"`
		Max int64 `json:"max"`
	} `json:"created_at"`
}

type ExportListResponse struct {
	baseResponse

	Data struct {
		Count int `json:"count"`
		List  []struct {
			Id            int64  `json:"id"`
			FileName      string `json:"name"`
			FileType      uint8  `json:"type"`
			FailReason    string `json:"fail_reason"`
			Process       uint8  `json:"process"`
			OperationName string `json:"operation_name"`
			OperationId   int64  `json:"operation_id"`
			Path          string `json:"path"`
			Status        uint8  `json:"status"`
			CreatedAt     int64  `json:"created_at"`
		} `json:"list"`
	} `json:"data"`
}

type RetryRequest struct {
	FileId int64 `json:"file_id"`
}

type RetryResponse struct {
	sdk.ResponseBase
}

type ProcessRequest struct {
}

type ProcessResponse struct {
	sdk.ResponseBase
}

type DownloadRequest struct {
	FileId       int64 `json:"file_id"`
	DownloadType uint8 `json:"download_type"`
}

type DownloadResponse struct {
	baseResponse
	Data struct {
		FileExt  string `json:"fileExt"`
		FileName string `json:"filename"`
		Host     string `json:"host"`
		Path     string `json:"path"`
	} `json:"data"`
}

type DeleteRequest struct {
	FileId int64 `json:"file_id" binding:"required" desc:"文件Id"`
}

type DeleteResponse struct {
	baseResponse
}

type ExportDetailResponse struct {
	baseResponse
	Data struct {
		Id            int64  `json:"id"`
		FileName      string `json:"name"`
		FileType      uint8  `json:"type"`
		FailReason    string `json:"fail_reason"`
		Process       uint8  `json:"process"`
		OperationName string `json:"operation_name"`
		Status        uint8  `json:"status"`
		CreatedAt     int64  `json:"created_at"`
		CompleteAt    int64  `json:"complete_at"`
		Path          string `json:"path" db:"path"`
		Host          string `json:"host"`
		FileExt       string `json:"file_ext"`
	} `json:"data"`
}

type DetailRequest struct {
	FileId int64 `json:"file_id" binding:"required" desc:"文件Id"`
}

type LatestExportDetailRequest struct {
	Label       string `json:"label" db:"label" validate:"required" required:"标签字段必须"`
	OperationId int64  `json:"operation_id"`
}

type OperationReq struct {
	OperationId int64 `json:"operation_id"`
}

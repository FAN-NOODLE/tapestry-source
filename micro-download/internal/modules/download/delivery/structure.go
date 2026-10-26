package delivery

const (
	// Redis进度Key
	SourceRedisProcess = "download-file-generate-process-%d"
	// TypeExcel 文件类型:excel
	TypeExcel = 1
	// TypeCsv 文件类型:csv
	TypeCsv = 2
	// TypeZip 文件类型:zip
	TypeZip = 3
	// SourceTypeData 资源类型: 数据
	SourceTypeData = 1
	// SourceTypeSql 资源类型: SQL
	SourceTypeSql = 2
	// SourceTypeRestful 资源类型: 接口
	SourceTypeRestful = 3
	// SourceTypePack 资源类型打包
	SourceTypePack = 4
	// StatusCreated 状态: 已创建
	StatusCreated = 1
	// StatusGenerating 状态: 生成中
	StatusGenerating = 2
	// StatusComplete 状态: 已生成
	StatusComplete = 3
	// StatusFailed 状态: 生成失败
	StatusFailed = 4
	// StatusInvalid 状态: 已失效
	StatusInvalid = 5
	// DownloadTypePath 下载类型: 前端下载 只返回path
	DownloadTypePath = 1
	// DownloadTypeBin 下载类型: 后端下载 返回二进制文件
	DownloadTypeBin = 2
	// PackResponsePath 打包数据返回类型 路径 TODO 未实现
	PackResponsePath = 1
	// PackResponseBin 打包数据返回类型 二进制文件
	PackResponseBin = 2
	// PackParamsRuleEach 打包类型参数请求规则: 遍历
	PackParamsRuleEach = 1
	// PackParamsRuleIncrement 打包类型参数请求规则: 自增
	PackParamsRuleIncrement = 2
	// PackParamsRuleDecrement 打包类型参数请求规则: 自减
	PackParamsRuleDecrement = 3
	// ExcelMergeTypeColumn excel单元格合并类型 列
	ExcelMergeTypeColumn = 1
	// ExcelMergeTypeLine excel单元格合并类型 行 TODO 未实现
	ExcelMergeTypeLine = 2
	// ExcelMergeRuleEqual excel单元格合并规则 相等
	ExcelMergeRuleEqual = "equal"
	// RestfulPaginationTypePage Restful分页类型: 页码
	RestfulPaginationTypePage = 1
	// RestfulPaginationTypeID Restful分页类型: ID
	RestfulPaginationTypeID = 2
)

type ExcelSourceDataCommonOption struct {
	SourceTitle          []string             `json:"source_title"`
	SourceTitleAlias     map[string]string    `json:"source_title_alias"`
	SourceTitleValue     []string             `json:"source_title_value"`
	SourceSuffix         [][]string           `json:"source_suffix"`
	SourceExcelCellMerge SourceExcelCellMerge `json:"source_excel_cell_merge"`
	SourceHeaderPrefix   [][]string           `json:"source_header_prefix"`
}

// SourceTypeDataData 数据类型
type SourceTypeDataData struct {
	ExcelSourceDataCommonOption
	Data [][]string
}

// SourceTypeSqlData sql类型
type SourceTypeSqlData struct {
	ExcelSourceDataCommonOption
	Count    int64  `json:"count"`
	Sql      string `json:"sql"`
	Database string `json:"database"`
}

// SourceTypeRestfulData 接口类型
type SourceTypeRestfulData struct {
	ExcelSourceDataCommonOption
	Path               string                 `json:"path"`
	Method             string                 `json:"method"`
	Headers            map[string]string      `json:"headers"`
	Params             map[string]interface{} `json:"params"`
	RequestPageField   string                 `json:"request_page_field"`
	RequestSizeField   string                 `json:"request_size_field"`
	ResponseCountField string                 `json:"count_field"`
	ResponseDataField  string                 `json:"data_field"`
	RequestFrom        []PageFromInfo         `json:"request_from"`
	SourceRestfulIDPaginationType
}

// SourceTypePackData 打包类型
type SourceTypePackData struct {
	ID      int64             `json:"id"`
	Path    string            `json:"path"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	SourceTypePackDataParams
}

// SourceTypePackDataParams 打包类型参数
type SourceTypePackDataParams struct {
	ParamsEach   map[string][]string                           `json:"params_each"`
	ParamsNumber map[string]SourceTypePackDataParamsNumberItem `json:"params_number"`
	ParamsRule   int8                                          `json:"params_rule"`
}

// SourceTypePackDataParamsNumberItem 打包参数item
type SourceTypePackDataParamsNumberItem struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
	Step  int   `json:"step"`
}

type CreateExport struct {
	ID                            int64                     `json:"id"`
	Filename                      string                    `json:"filename" validate:"required" required:"文件名称[filename]是必填项"`
	Type                          int8                      `json:"type"`
	EnterpriseID                  int64                     `json:"enterprise_id" validate:"required" required:"企业ID[enterprise_id]是必填项"`
	StoreID                       int64                     `json:"store_id" validate:"required" required:"门店ID[store_id]是必填项"`
	OperationID                   int64                     `json:"operation_id" validate:"required" required:"操作人ID[operation_id]是必填项"`
	OperationName                 string                    `json:"operation_name" validate:"required" required:"操作人名称[operation_name]是必填项"`
	SourceType                    int8                      `json:"source_type" validate:"required" required:"数据源类型[source_type]是必填项"`
	SourceTitle                   []string                  `json:"source_title"`
	SourceTitleAlias              map[string]string         `json:"source_title_alias"`
	SourceTitleValue              []string                  `json:"source_title_value"`
	SourceSuffix                  [][]string                `json:"source_suffix"`
	SourceExcelCellMerge          SourceExcelCellMerge      `json:"source_excel_cell_merge"`
	SourceData                    [][]string                `json:"source_data"`
	SourceSqlDatabase             string                    `json:"source_sql_database"`
	SourceSqlCount                int64                     `json:"source_sql_count"`
	SourceSql                     string                    `json:"source_sql"`
	SourceRestfulPath             string                    `json:"source_restful_path"`
	SourceRestfulMethod           string                    `json:"source_restful_method"`
	SourceRestfulHeaders          map[string]string         `json:"source_restful_headers"`
	SourceRestfulParams           map[string]interface{}    `json:"source_restful_params"`
	SourceRestfulCountField       string                    `json:"source_restful_count_field"`
	SourceRestfulDataField        string                    `json:"source_restful_data_field"`
	SourceRestfulRequestPageField string                    `json:"source_restful_request_page_field"`
	SourceRestfulRequestSizeField string                    `json:"source_restful_request_size_field"`
	SourceTypePackPath            string                    `json:"source_pack_path"`
	SourceTypePackMethod          string                    `json:"source_type_pack_method"`
	SourceTypePackHeaders         map[string]string         `json:"source_type_pack_headers"`
	SourceTypePackParams          *SourceTypePackDataParams `json:"source_type_pack_params"`
	SourceHeaderPrefix            [][]string                `json:"source_header_prefix"`
	Label                         string                    `json:"label"`

	// 关于id分段的一些可传入的值
	SourceRestfulIDPaginationType
}

type SourceRestfulIDPaginationType struct {
	SourceRestfulPaginationType      int8              `json:"source_restful_pagination_type"`       // =2代表是id分段，其他就是默认值
	SourceRestfulReceiveSelectSwitch bool              `json:"source_restful_receive_select_switch"` // true 代表我们数据源在查询时是可以精简查询字段的
	SourceRestfulReceiveSelectField  string            `json:"source_restful_receive_select_field"`  // 我们将要接收的字段，我们可以忽略这个，优化一个点
	SourceRestfulSortField           string            `json:"source_restful_sort_field"`            // 这是我们的分类id
	SourceRestfulSortKeyAs           string            `json:"source_restful_sort_key_as"`           // 我们的排序字段在结构体中是什么
	SourceRestfulSortKey             string            `json:"source_restful_sort_key"`
	SourceRestfulSortTypeField       string            `json:"source_restful_sort_type_field"` // 这个是我们的分类方式字段
	SourceRestfulSortType            string            `json:"source_restful_sort_type"`
	SourceRestfulScrollPageInfoField string            `json:"source_restful_scroll_page_info_field"` // 这个代表我们当前在分页，是一个bool，我们应该返回从当前游标位置走一个size大小的距离，返回这个id即可
	SourceRestfulSelectInfo          []string          `json:"source_restful_select_info"`            // 这个是我们可以提前约定一个select，如果不约定我们会按照title来走
	SourceRestfulMinIDField          string            `json:"source_restful_min_id_field"`
	SourceRestfulMaxIDField          string            `json:"source_restful_max_id_field"`
	SourceRestfulExt                 map[string]string `json:"source_restful_ext"`
}

type SourceExcelCellMerge struct {
	Type       uint8                  `json:"type"`
	Rule       string                 `json:"rule"`
	Title      []string               `json:"title"`
	Premise    []string               `json:"premise"`
	PremiseMap map[string]interface{} `json:"premisemap"`
}

type CreatedExport struct {
	FileID int64 `json:"file_id"`
}

type ExportListReq struct {
	Filename      string  `json:"filename"`
	EnterpriseID  int64   `json:"enterprise_id" validate:"required" required:"企业ID[enterprise_id]是必填项"`
	StoreIDS      []int64 `json:"store_ids"`
	Status        int8    `json:"status"`
	OperationIDS  []int64 `json:"operation_ids"`
	OperationId   int64   `json:"operation_id"`
	OperationName string  `json:"operation_name"`
	CreatedAt     struct {
		Min int64 `json:"min"`
		Max int64 `json:"max"`
	} `json:"created_at"`
	Page     int64 `json:"page"`
	PageSize int64 `json:"page_size"`
}

type ExportList struct {
	ID            int64  `json:"id" db:"id"`
	Filename      string `json:"name" db:"name"`
	Type          int8   `json:"type" db:"type"`
	OperationName string `json:"operation_name" db:"operation_name"`
	OperationId   int64  `json:"operation_id" db:"operation_id"`
	Path          string `json:"path" db:"path"`
	Status        int8   `json:"status" db:"status"`
	CreatedAt     int64  `json:"created_at" db:"created_at"`
	FailedReason  string `json:"failed_reason" db:"failed_reason"`
	Process       int    `json:"process"`
	CompleteAt    int64  `json:"complete_at" db:"complete_at"`
	Host          string `json:"host"`
	FileExt       string `json:"file_ext"`
}

type ExportListRes struct {
	Count int64        `json:"count"`
	List  []ExportList `json:"list"`
}

type FileID struct {
	FileID int64 `json:"file_id" db:"id" validate:"required" required:"主键ID必须"`
}

type Label struct {
	Label       string `json:"label" db:"label" validate:"required" required:"标签字段必须"`
	OperationId int64  `json:"operation_id" db:"operation_id"`
}

type PageFromInfo struct {
	MinID int64
	MaxID int64
	Page  int // 合理情况下是不包含page和size的
	Size  int // 这俩东西是我们的sortKey如果不是唯一id，并且一个id还能跑出比我们默认pageSize大小的时候会用到
}

type SourceOverview struct {
	ID         int64
	Count      int64
	TitleAlias map[string]string
	Title      []string
	Value      []string
	PageNum    int64
	PageFrom   []PageFromInfo
	Suffix     [][]string
	HeadPrefix [][]string
	CellMerge  SourceExcelCellMerge
}

type ProcessRes struct {
	Process int `json:"process"`
}

type DownloadReq struct {
	FileID       int64 `json:"file_id" validate:"required" required:"主键ID必须"`
	DownloadType int8  `json:"download_type" validate:"required" required:"下载类型必须"`
}

type DownloadData struct {
	DownloadType int8   `json:"download_type"`
	ContentType  string `json:"content_type"`
	Filename     string `json:"filename"`
	FileExt      string `json:"file_ext"`
	FileData     []byte `json:"bin"`
	Path         string `json:"path"`
}

type OperationReq struct {
	OperationId int64 `json:"operation_id" db:"operation_id"`
}

package download

type ExportListReq struct {
	Page          int     `json:"page"`
	PageSize      int     `json:"page_size"`
	Status        uint8   `json:"status"`
	Filename      string  `json:"filename"`
	OperationIds  []int64 `json:"operation_ids"`
	OperationName string  `json:"operation_name"`
	CreatedAt     struct {
		Min int64 `json:"min"`
		Max int64 `json:"max"`
	} `json:"created_at"`
}

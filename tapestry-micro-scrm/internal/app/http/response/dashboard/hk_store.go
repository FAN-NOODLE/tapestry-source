package dashboardResp

type HkStoreDetail struct {
	ID                 int64  `json:"id"`
	StoreCode          string `json:"store_code"`
	WechatDepartmentID uint   `json:"wechat_department_id"`
	StoreName          string `json:"store_name"`
	Area               string `json:"area"`
	StoreType          int32  `json:"store_type"`
	InventoryName      string `json:"inventory_name"`
	Address            string `json:"address"`
	City               string `json:"city"`
	Phone              string `json:"phone"`
	Email              string `json:"email"`
	ClosedAt           int64  `json:"closed_at"`
	CreatedAt          int64  `json:"created_at"`
	UpdatedAt          int64  `json:"updated_at"`
}

type HkStoreSyncResult struct {
	CreatedCount int `json:"created_count"`
	SkippedCount int `json:"skipped_count"`
	ClosedCount  int `json:"closed_count"`
}

type HkStoreLastSyncTime struct {
	LastSyncTime int64 `json:"last_sync_time"`
}

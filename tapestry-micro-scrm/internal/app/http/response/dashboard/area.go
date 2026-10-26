package dashboardResp

type AreaList struct {
	ID        int    `json:"id"`
	Type      int    `json:"type" `
	Name      string `json:"name"`
	ParentID  int    `json:"parent_id"`
	ThirdID   int    `json:"third_id"`
	ClosedAt  int64  `json:"closed_at"`
	StoreCode string `json:"store_code"`
	StoreName string `json:"store_name"`
	PathName  string `json:"path_name"`
	UpdatedAt int64  `json:"updated_at"`
	//获取子区域
	Children []AreaList `json:"children"`
}

type AreaInfo struct {
	ID               int64  `json:"id"`
	Type             int    `json:"type" `
	Name             string `json:"name"`
	ChildrenThirdIds []int  `json:"children_third_ids"`
}

type FileRow struct {
	CustomGroup        string `json:"custom_group"`
	CustomGroupElement string `json:"custom_group_element" `
	Location           string `json:"location"`
	RetailPOSNumber    string `json:"retail_pos_number"`
}

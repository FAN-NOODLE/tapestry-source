package delivery

type SignReceiveInfo struct {
	ID        int64 `json:"id"`
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
	Type      int   `json:"type"`
	Page      Page  `json:"page"`
}

type SignDetailInfo struct {
	ID        int64  `json:"id"`         //签到活动的id
	UnionID   string `json:"union_id"`   //
	StartTime int64  `json:"start_time"` //开始时间
	EndTime   int64  `json:"end_time"`   // 结束时间
	Page      Page   `json:"page"`       //分页
}

type SignInList struct {
	UnionID      string `json:"union_id"`
	Count        int    `json:"count"`
	LastSignTime int64  `json:"last_sign_time"`
	Name         string `json:"name"`
	Phone        string `json:"phone"`
}

type NewResultList struct {
	Count      int64        `json:"count"`
	SignInList []SignInList `json:"sign_in_list"`
}

type SignInInfoDetail struct {
	Name        string `json:"wx_name"`
	Phone       string `json:"wx_phone"`
	Day         int    `json:"day"`
	SignTime    int64  `json:"sign_time"`
	SignContent string `json:"sign_content"`
	SignPic     string `json:"sign_pic"`
	Shr         string `json:"shr"`
	Mobile      string `json:"mobile"`
	Address     string `json:"address"`
}
type Page struct {
	PageSize int `json:"page_size"`
	PageNum  int `json:"page_num"`
}

package activityQRCodeResp

type ActivityInfoResp struct {
	ActivityBaseResp
	TagList         []Info                `json:"tag_list"`
	ScopeMemberList []Info                `json:"scope_member_list"`
	ScopeAreaList   []Info                `json:"scope_area_list"`
	CategoryList    []CategoryInfo        `json:"category_list"`
	WelcomeList     []ActivityWelcomeInfo `json:"welcome_list"`
}
type ActivityBaseResp struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	AliasName      string `json:"alias_name"`
	BannerUrl      string `json:"banner_url"`
	TimeType       int    `json:"time_type"`
	StartTime      int    `json:"start_time"`
	EndTime        int    `json:"end_time"`
	SkipVerify     int    `json:"skip_verify"`
	AddNameTag     int    `json:"add_name_tag"`
	AddCategoryTag int    `json:"add_category_tag"`
}

type ActivityListResp struct {
	ActivityBaseResp
	CreatedAt     int64    `json:"created_at"`
	TagNames      []string `json:"tag_names"`
	CategoryNames []string `json:"category_names"`
	HasWelcome    int      `json:"has_welcome"`
}
type Info struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type PointInfo struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	CodeUrl  string `json:"code_url"`
	LinkUrl  string `json:"link_url"`
	QueryStr string `json:"query_str"`
}

type CategoryInfo struct {
	ID        int64       `json:"id"`
	Name      string      `json:"name"`
	LinkUrl   string      `json:"link_url"`
	CodeUrl   string      `json:"code_url"`
	QueryStr  string      `json:"query_str"`
	PointList []PointInfo `json:"point_list"`
}
type ActivityWelcomeInfo struct {
	Type    int                      `json:"type" commit:"0 通用 1 会员 2非会员"`
	Content []map[string]interface{} `json:"content" db:"content" binding:"required" desc:"内容"`
}

type WelcomeInfo struct {
	ID           int64                    `json:"id"`
	CategoryId   int64                    `json:"category_id"`
	CategoryName string                   `json:"category_name"`
	Type         int                      `json:"type" commit:"0 通用 1 会员 2非会员"`
	Content      []map[string]interface{} `json:"content" db:"content" binding:"required" desc:"内容"`
}

type WelcomeListResp struct {
	ID           int64                    `json:"id"`
	CategoryName string                   `json:"category_name"`
	CreatedAt    int64                    `json:"created_at"`
	Content      []map[string]interface{} `json:"content" db:"content" binding:"required" desc:"内容"`
}

type QrCodeResp struct {
	ErrorCode        int              `json:"error_code"`
	ActivityName     string           `json:"activity_name"`
	BannerUrl        string           `json:"banner_url"`
	QrCodeMemberInfo QrCodeMemberInfo `json:"qr_code_member_info"`
	StoreList        []StoreInfo      `json:"store_list"`
}
type QrCodeMemberInfo struct {
	MemberId  int64  `json:"member_id"`
	UserId    string `json:"user_id"`
	QrCode    string `json:"qr_code"`
	StoreCode string `json:"store_code"`
	StoreName string `json:"store_name"`
	ThirdId   int    `json:"third_id"`
	SaName    string `json:"sa_name"`
	Avatar    string `json:"avatar"`
}

type StoreListResp struct {
	ErrorCode int         `json:"error_code"`
	StoreList []StoreInfo `json:"store_list"`
}

type StoreInfo struct {
	ThirdId      int    `json:"third_id"`
	StoreCode    string `json:"store_code"`
	StoreName    string `json:"store_name"`
	StoreAddress string `json:"store_address"`
	City         string `json:"city"`
	CityPinYin   string `json:"city_pin_yin"`
	Distance     int    `json:"distance"`
}

type StatisticsResp struct {
	TotalStatistics TotalStatistics   `json:"total_statistics"`
	List            []PointStatistics `json:"list"`
	Total           int               `json:"total"`
	Page            int               `json:"page"`
	Size            int               `json:"size"`
}

type TotalStatistics struct {
	State            string `json:"state"`
	AddRelationTotal int    `json:"add_relation_total"`
	AddCustomerTotal int    `json:"add_customer_total"`
	DelRelationTotal int    `json:"del_relation_total"`
	AddVipTotal      int    `json:"add_vip_total"`
}
type PointStatistics struct {
	ActivityName    string `json:"activity_name"`
	CategoryName    string `json:"category_name"`
	PointName       string `json:"point_name"`
	TotalStatistics `json:"total_statistics"`
}

const ActivityErrNotExists = 1
const ActivityErrNotStart = 2
const ActivityErrExpired = 3
const ActivityErrNoSa = 4
const ActivityErrWeComErr = 5

package storeVisitResp

type SendEmail struct {
	Emails    []string `json:"emails"`
	StoreName string   `json:"store_name"`
	TimeStr   string   `json:"time_str"`
}

type DeptInfo struct {
	DeptName   string       `json:"dept_name"`
	MemberList []MemberInfo `json:"member_list"`
}

type EmailContact struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

type MemberInfo struct {
	MemberId     int64  `json:"member_id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	DeptName     string `json:"dept_name"`
	StoreName    string `json:"store_name"`
	UserId       string `json:"user_id"`
	Avatar       string `json:"avatar"`
	IsMM         int    `json:"is_mm"`
	PositionAbbr string `json:"position_abbr"`
	ThirdId      int    `json:"third_id"`
}

type LoginUserInfo struct {
	MemberInfo
	ThirdId             int    `json:"third_id"`
	StoreCode           string `json:"store_code"`
	Avatar              string `json:"avatar"`
	Position            string `json:"position"`
	IsStoreVisitManager bool   `json:"is_store_visit_manager"`
	IsStoreManager      bool   `json:"is_store_manager"`
}

type StoreVisitMember struct {
	ID        int64  `json:"id"`
	Date      string `json:"date"`
	TimeSpan  string `json:"time_span"`
	ThirdId   int    `json:"third_id"`
	Remark    string `json:"remark"`
	Status    int    `json:"status"`
	UpdatedAt int64  `json:"updated_at"`
}

type TotalCountByStore struct {
	TotalCount int64 `json:"total_count"`
	ThirdId    int   `json:"third_id"`
}

type StoreVisitInfo struct {
	ID            int64        `json:"id"`
	Date          string       `json:"date"`
	TimeSpan      string       `json:"time_span"`
	ThirdId       int          `json:"third_id"`
	StoreName     string       `json:"store_name"`
	StoreCode     string       `json:"store_code"`
	Remark        string       `json:"remark"`
	Status        int          `json:"status"`
	Creator       MemberInfo   `json:"creator"`
	StoreManager  []MemberInfo `json:"store_manager"`
	MMManager     []MemberInfo `json:"mm_manager"`
	VisitorMember []MemberInfo `json:"visitor_member"`
	CopyMember    []MemberInfo `json:"copy_member"`
	ToStoreTime   int          `json:"to_store_time"`
	VisitorIsShow int          `json:"visitor_is_show"`
	StoreIsShow   int          `json:"store_is_show"`
	Overview      string       `json:"overview"`
	Content       string       `json:"content"`
	PdfUrl        string       `json:"pdf_url"`
	IsReport      int          `json:"is_report"`
	UpdatedAt     int64        `json:"updated_at"`
}

type ShowRecordMediaInfo struct {
	ID      int64 `json:"id"`
	MediaID int64 `json:"media_id"`
}

type AdminStoreVisitInfo struct {
	ID        int64  `json:"id"`
	StoreName string `json:"store_name"`
	StoreCode string `json:"store_code"`
	Status    int    `json:"status"`
	Date      string `json:"date"`
	TimeSpan  string `json:"time_span"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type VisitCommentInfo struct {
	ID              int64       `json:"id"`
	VisitID         int64       `json:"visit_id"`
	MemberInfo      MemberInfo  `json:"member_info"`
	Content         string      `json:"content"`
	ParentID        int64       `json:"parent_id"`
	ReplyMemberInfo MemberInfo  `json:"reply_member_info"`
	IsDel           int         `json:"is_del"`
	MediaList       []MediaInfo `json:"media_list"`
	CreatedAt       int64       `json:"created_at"`
}

type MediaInfo struct {
	MediaType     int    `json:"media_type"`
	Url           string `json:"url"`
	VideoImageUrl string `json:"video_image_url"`
	VideoSeconds  int    `json:"video_second"`
}
type VisitCommentAllInfo struct {
	VisitCommentInfo
	Children []VisitCommentInfo `json:"children"`
}

type AdminDashboardInfo struct {
	Name       string `json:"name" comment:"这里门店汇总，指的是门店名称，渠道汇总指的是渠道名称，大区汇总指的是大区名称，小区汇总指的是小区名称"`
	ParentName string `json:"parent_name" comment:"小区汇总，展示大区名称"`
	StoreCode  string `json:"store_code"`
	DeptName   string `json:"dept_name"`
	VisitNum   int64  `json:"visit_num"`
	VisitRate  string `json:"visit_rate"`
}

type AdminDashboardStore struct {
	MemberId  int64 `json:"member_id"`
	ThirdId   int   `json:"third_id"`
	VisitNum  int64 `json:"visit_num"`
	FinishNum int64 `json:"finish_num"`
}

type DashboardData struct {
	Name       string `json:"name"`
	TotalCount int64  `json:"total_count"`
}

type DashboardTree struct {
	Name     string          `json:"name"` // 类型名称，如"奥莱"或"零售"
	Level    int             `json:"level"`
	Children []*ProvinceNode `json:"children"` // 省份节点
	DataList []DashboardData `json:"dataList"`
}

type ProvinceNode struct {
	Name       string      `json:"name"`
	Level      int         `json:"level"`
	ThirdIds   []int       `json:"third_ids"` //所有该省份下的wechat_department_id
	TotalCount int64       `json:"total_count"`
	Children   []*CityNode `json:"children"` // 城市节点

}

type CityNode struct {
	Name       string       `json:"name"`
	Level      int          `json:"level"`
	ThirdIds   []int        `json:"third_ids"` //所有该城市下的wechat_department_id
	TotalCount int64        `json:"total_count"`
	Children   []*StoreNode `json:"children"` // 门店节点
}
type StoreNode struct {
	Name       string `json:"name"`
	StoreCode  string `json:"store_code"`
	Level      int    `json:"level"`
	ThirdId    int    `json:"third_id"`
	TotalCount int64  `json:"total_count"`
}

type ThirdIdAreaInfo struct {
	Source    string `json:"source"`     // 类型名称，如"奥莱"或"零售"
	Region    string `json:"region"`     // 大区
	SubRegion string `json:"sub_region"` //分区
}

type AreaInfo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type TotalCountByStoreAndMember struct {
	TotalCount int64 `json:"total_count"`
	MemberId   int64 `json:"member_id"`
	ThirdId    int   `json:"third_id"`
}

type DashboardExportKey struct {
	DeptName string `json:"dept_name"`
	ThirdId  int    `json:"third_id"`
}

type DashboardExportInfo struct {
	Source        string `json:"source"`
	Region        string `json:"region"`
	SubRegion     string `json:"sub_region"`
	City          string `json:"city"`
	Province      string `json:"province"`
	StoreCode     string `json:"store_code"`
	StoreName     string `json:"store_name"`
	VisitDeptName string `json:"visit_dept_name"`
	TotalCount    int64  `json:"total_count"`
}

type PdfReportInfo struct {
	Date            string          `json:"date"`
	TimeSpan        string          `json:"time_span"`
	StoreName       string          `json:"store_name"`
	CreatorName     string          `json:"creator_name"`
	CreatorDeptName string          `json:"creator_dept_name"`
	Overview        string          `json:"overview"`
	Content         string          `json:"content"`
	ShowRecordList  []PdfShowRecord `json:"show_record_list"`
	UpdatedAt       int64           `json:"updated_at"`
}

type PdfShowRecord struct {
	TypeName  string     `json:"type_name"`
	Content   string     `json:"content"`
	MediaList []PdfMedia `json:"media_list"`
}
type PdfMedia struct {
	MediaType     int    `json:"media_type"`
	URL           string `json:"url"`
	ThumbUrl      string `json:"thumb_url"`
	VideoImageUrl string `json:"video_image_url"`
}

type ConfigInfoResp struct {
	Key          string `json:"Key"`
	Value        string `json:"value"`
	MemberConfig string `json:"member_config"`
}

type StoreMMResp struct {
	ThirdId     int     `json:"third_id"`
	MMMemberIds []int64 `json:"mm_member_ids"`
}

const (
	LevelSource   = iota + 1 //渠道
	LevelProvince            //省
	LevelCity                //市
	LevelStore               //门店
)

package source

import (
	"micro-scrm/internal/app/models"
)

const (
	// TypeText 文本类型
	TypeText = iota + 1
	// TypeImg 图片类型
	TypeImg
	// TypeLink 链接类型
	TypeLink
	// TypeMiniProgram 小程序类型
	TypeMiniProgram
	// TypeVideo 视频类型
	TypeVideo
	// TypeFile 文件类型
	TypeFile
	// TypeFriend 朋友圈
	TypeFriend
	// TypeProduct 官网商品
	TypeProduct
)
const (
	// StatusEnable 启用0
	StatusEnable = iota
	// StatusDisable 禁用1
	StatusDisable
)
const (
	// Deal 压缩的状态
	Deal = iota
	DealING
	DealSuccess
	DealFail
)

const (
	TopNo = iota
	// 置顶
	TopYes
)

type SourceModel struct {
	models.MicroBaseModel
	GroupID         int64 `json:"group_id" gorm:"column:group_id;type:bigint(20);not null;index:group_id_index;default:0;comment:分组ID 0 未分组"`
	Type            int32 `json:"type" gorm:"column:type;type:tinyint(3);not null;default:0;comment:类型(1 文本 2 图片 3 链接 4 小程序 5 视频 6 文件)"`
	Status          int32 `json:"status" gorm:"column:status;type:tinyint(3);not null;default:0;comment:启用0；禁用1"`
	AdminID         int64 `json:"admin_id" gorm:"column:admin_id;type:bigint(20);not null;index:admin_id_index;default:0;comment:操作人ID"`
	StartTime       int   `json:"start_time" gorm:"column:start_time;type:int(11);not null;default:0;comment:开始时间"`
	EndTime         int   `json:"end_time" gorm:"column:end_time;type:int(11);not null;default:0;comment:结束时间"`
	MemberID        int64 `json:"member_id" gorm:"column:member_id;type:bigint(20);not null;default:0;comment:导购的id，member_id"`
	IsDeal          int32 `json:"is_deal" gorm:"column:is_deal;type:tinyint(3);not null;default:0;comment:压缩状态,0为不压缩,1为压缩中,2压缩成功,3压缩失败"`
	CompressID      int64 `json:"compress_id" gorm:"column:compress_id;type:bigint(20);not null;default:0;comment:压缩素材集合id"`
	MergeID         int64 `json:"merge_id" gorm:"column:merge_id;type:bigint(20);not null;default:0;comment:合并组id"`
	Sort            int32 `json:"sort" gorm:"column:sort;type:tinyint(1);not null;default:0;comment:排序"`
	DownloadTimes   int   `json:"download_times" gorm:"column:download_times;type:int(10);not null;default:0;comment:下载次数"`
	ForwardingTimes int   `json:"forwarding_times" gorm:"column:forwarding_times;type:int(10);not null;default:0;comment:转发次数"`
	CopyTimes       int   `json:"copy_times" gorm:"column:copy_times;type:int(10);not null;default:0;comment:复制次数"`
	CollectNum      int   `json:"collect_num" gorm:"column:collect_num;type:int(10);not null;default:0;comment:收藏数"`
	IsTop           int   `json:"is_top" gorm:"column:is_top;type:tinyint(1);not null;default:0;comment:是否上架"`
}

// TableName 素材主表
func (SourceModel) TableName() string {
	return "source"
}

type NewSourceList struct {
	IdCount   int64 `json:"id_count"`
	ThirdId   int   `json:"third_id"`
	StartTime int64 `json:"start_time"`
}

type NewSourceCountDetail struct {
	IdCount      int64  `json:"id_count"`
	DepartmentId int64  `json:"department_id"`
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	Userid       string `json:"userid"`
	CreatedAt    int64  `json:"created_at"`
}

type NewSourceCountNum struct {
	CountNum int64 `json:"count_num"`
	Type     int   `json:"type"`
}

type SourceCollectNumInfo struct {
	Id         int64 `json:"id"`
	MergeId    int64 `json:"merge_id"`
	CollectNum int64 `json:"collect_num"`
}

type SourceProductCollectNumInfo struct {
	ProductCode string `json:"product_code"`
	CollectNum  int64  `json:"collect_num"`
}

package microScrmModel

type MemberLearnRankModel struct {
	ID           int    `gorm:"primaryKey" json:"id"`                        // 主键
	MemberID     int64  `gorm:"column:member_id" json:"member_id"`           // 导购 member 表 id
	DepartmentID int64  `gorm:"column:department_id" json:"department_id"`   // departments 表主键
	FileType     string `gorm:"column:file_type;size:32" json:"file_type"`   // 排行类型
	PMonth       int    `gorm:"column:p_month" json:"p_month"`               // P 几财月的
	FyYear       int    `gorm:"column:fy_year" json:"fy_year"`               // 财年
	UserID       string `gorm:"column:user_id;size:32" json:"user_id"`       // 导购编码
	StoreCode    string `gorm:"column:store_code;size:32" json:"store_code"` // 门店编码
	Value        string `gorm:"column:value;size:64" json:"value"`           // 具体数值
	Rank         int    `gorm:"column:rank" json:"rank"`                     // 排行
	Remark       string `gorm:"column:remark;size:255" json:"remark"`        // 备注
	UsedAt       int64  `gorm:"column:used_at" json:"used_at"`               // 消费时间
	CreatedAt    int64  `gorm:"column:created_at" json:"created_at"`         // 创建时间
	UpdatedAt    int64  `gorm:"column:updated_at" json:"updated_at"`         // 修改时间
	DeletedAt    int64  `gorm:"column:deleted_at" json:"deleted_at"`         // 删除时间
}

// TableName returns the table name of the MemberLearnRecord model
func (MemberLearnRankModel) TableName() string {
	return "member_learn_rank"
}

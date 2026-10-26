package micro_scrm_enterprise

import "micro-scrm/internal/orm/model"

type BaKanbanDictionaryModel struct {
	model.BranchModel
	BaName                string `gorm:"column:ba_name;NOT NULL" json:"ba_name"`                        // 名称
	BaMobile              string `gorm:"column:ba_mobile" json:"ba_mobile"`                             // 手机号码
	BaLevel               int    `gorm:"column:ba_level;NOT NULL" json:"ba_level"`                      // 级别 1 ba 、2 门店经理、3分区经理、4大区经理
	CounterName           string `gorm:"column:counter_name" json:"counter_name"`                       // 柜台名称
	CounterNo             string `gorm:"column:counter_no" json:"counter_no"`                           // 柜台编号
	ManagerName           string `gorm:"column:manager_name" json:"manager_name"`                       // 店长名称
	ManagerMobile         string `gorm:"column:manager_mobile" json:"manager_mobile"`                   // 店长电话
	PartitionLeaderName   string `gorm:"column:partition_leader_name" json:"partition_leader_name"`     // 分区领导姓名
	PartitionLeaderMobile string `gorm:"column:partition_leader_mobile" json:"partition_leader_mobile"` // 分区领导电话
	RegionalLeaderName    string `gorm:"column:regional_leader_name" json:"regional_leader_name"`       // 大区领导
	RegionalLeaderMobie   string `gorm:"column:regional_leader_mobie" json:"regional_leader_mobie"`     // 大区领导电话
	RegionalArea          string `gorm:"column:regional_area" json:"regional_area"`                     // 大区区域
}

func (BaKanbanDictionaryModel) TableName() string {
	return "ba_kanban_dictionary"
}

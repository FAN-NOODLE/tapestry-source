package excitationReq

type MemberDataListReq struct {
	EnterpriseId          int64  `json:"enterprise_id" comment:"企业ID"`
	Page                  int    `json:"page" binding:"min=1" comment:"页码"`
	Size                  int    `json:"size" binding:"min=10,max=1000" comment:"分页数"`
	Total                 int    `json:"total"  comment:"总数"`
	MemberName            string `json:"member_name"`
	MemberUserid          string `json:"member_userid"`
	MemberID              int    `json:"member_id"`
	MemberThirdID         int    `json:"member_third_id"  comment:"third_id"`
	Level                 string `json:"level"  comment:"等级名称"`
	PointUpdatedStartTime int    `json:"point_updated_start_time"  comment:"积分更新时间"`
	PointUpdatedEndTime   int    `json:"point_updated_end_time" comment:"积分更新时间"`
	LevelUpdatedStartTime int    `json:"level_updated_start_time"  comment:"等级更新时间"`
	LevelUpdatedEndTime   int    `json:"level_updated_end_time" comment:"等级更新时间"`
	MedalUpdatedStartTime int    `json:"medal_updated_start_time"  comment:"等级更新时间"`
	MedalUpdatedEndTime   int    `json:"medal_updated_end_time" comment:"等级更新时间"`
	PointValue            string `json:"point_value"  comment:"积分数量排序"`
	MedalNum              string `json:"medal_num"  comment:"奖章数量排序"`
	Coins                 string `json:"coins"  comment:"金币排序"`
	SpendCoins            string `json:"spend_coins"  comment:"累计消耗金币排序"`
	TotalCoins            string `json:"total_coins"  comment:"累计总金币排序"`
}

type MemberDataInfoResp struct {
	Id                int64   `json:"id"`
	MemberID          int64   `json:"member_id"`
	MemberUserid      string  `json:"member_userid"`
	MemberName        string  `json:"member_name"`
	ThirdName         string  `json:"third_name"`
	ThirdId           int     `json:"third_id"`
	PointValue        string  `json:"point_value" gorm:"comment:积分(整数)"`
	PointUpdatedAt    int     `json:"point_updated_at" gorm:"comment:积分更新时间"`
	Coins             float64 `json:"coins" gorm:"comment:金币"`
	TotalCoins        float64 `json:"total_coins" gorm:"comment:总获得金币"`
	SpendCoins        float64 `json:"spend_coins" gorm:"comment:总花费金币"`
	PointUpdatedAtStr string  `json:"point_updated_at_str" gorm:"comment:积分更新时间"`
	Level             string  `json:"level" gorm:"comment:等级名称"`
	LevelID           int64   `json:"level_id" gorm:"comment:level_rule表Id"`
	LevelUpdatedAt    int     `json:"level_updated_at" gorm:"comment:等级更新时间"`
	LevelUpdatedAtStr string  `json:"level_updated_at_str" gorm:"comment:更新时间"`
	MedalNum          int     `json:"medal_num" gorm:"comment:奖章数量(整数)"`
	MedalUpdatedAt    int     `json:"medal_updated_at" gorm:"comment:奖章更新时间"`
	MedalUpdatedAtStr string  `json:"medal_updated_at_str" gorm:"comment:更新时间"`
	MaxLevel          string  `json:"max_level"`
	MaxLevelTime      int64   `json:"max_level_time"`
	MaxLevelTimeStr   string  `json:"max_level_time_str"`
}

type MemberDataListResp struct {
	Total int `json:"total"`
	List  []MemberDataInfoResp
}

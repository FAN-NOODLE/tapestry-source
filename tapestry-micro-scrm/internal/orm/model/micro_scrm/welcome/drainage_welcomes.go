package welcome

import (
	"micro-scrm/internal/orm/model"
)

const (
	TypeSelfHoo  = 10 // 个性欢迎语
	CHANNLEID_KD = 0
	ISDEFAULT_KD = 0
)

type DrainageWelcomesModel struct {
	model.MicroScrmBaseModel

	Content          string `json:"content" gorm:"column:content;type:text;not null;comment:内容"`
	IsDefault        int    `json:"is_default" gorm:"column:is_default;type:int(1);not null;default:0;comment:是否是默认欢迎语 1是 0不是"`
	ChannelID        int64  `json:"channel_id" gorm:"column:channel_id;type:bigint(20);not null;default:0;comment:渠道id,标识企业微信state"`
	Type             uint32 `json:"type" gorm:"column:type;type:tinyint(1) unsigned;not null;default:1;comment:欢迎语类型：1-系统欢迎语，2-裂变活动欢迎语，10-个性欢迎语"`
	ConditionContent string `json:"condition_content" gorm:"column:condition_content;type:text;not null;comment:条件集合，json格式"`
}

func (m DrainageWelcomesModel) TableName() string {
	return "drainage_welcomes"
}
